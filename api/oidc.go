package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/auth"
	"github.com/gotify/server/v2/config"
	"github.com/gotify/server/v2/database"
	"github.com/gotify/server/v2/decaymap"
	"github.com/gotify/server/v2/model"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func NewOIDC(conf *config.Configuration, db *database.GormDatabase, userChangeNotifier *UserChangeNotifier) *OIDCAPI {
	scopes := conf.OIDC.Scopes
	if len(scopes) == 0 {
		scopes = []string{"openid", "profile", "email"}
	}

	cookieKey := make([]byte, 32)
	if _, err := rand.Read(cookieKey); err != nil {
		log.Fatalf("failed to generate OIDC cookie key: %v", err)
	}
	cookieHandlerOpt := []httphelper.CookieHandlerOpt{}
	if !conf.Server.SecureCookie {
		cookieHandlerOpt = append(cookieHandlerOpt, httphelper.WithUnsecure())
	}
	cookieHandler := httphelper.NewCookieHandler(cookieKey, cookieKey, cookieHandlerOpt...)

	opts := []rp.Option{rp.WithCookieHandler(cookieHandler), rp.WithPKCE(cookieHandler)}

	provider, err := rp.NewRelyingPartyOIDC(
		context.Background(),
		conf.OIDC.Issuer,
		conf.OIDC.ClientID,
		conf.OIDC.ClientSecret,
		conf.OIDC.RedirectURL,
		scopes,
		opts...,
	)
	if err != nil {
		log.Fatalf("failed to initialize OIDC provider: %v", err)
	}

	return &OIDCAPI{
		DB:                 db,
		Provider:           provider,
		UserChangeNotifier: userChangeNotifier,
		UsernameClaim:      conf.OIDC.UsernameClaim,
		PasswordStrength:   conf.PassStrength,
		SecureCookie:       conf.Server.SecureCookie,
		AutoRegister:       conf.OIDC.AutoRegister,
		pendingSessions:    decaymap.NewDecayMap[string, *pendingOIDCSession](time.Now(), pendingSessionMaxAge),
	}
}

const pendingSessionMaxAge = 10 * time.Minute

type pendingOIDCSession struct {
	RedirectURI string
	ClientName  string
	CreatedAt   time.Time
	Elevate     *pendingElevation
}

type pendingElevation struct {
	ClientID        uint `form:"id" binding:"required"`
	DurationSeconds int  `form:"durationSeconds" binding:"required"`
}

// OIDCAPI provides handlers for OIDC authentication.
type OIDCAPI struct {
	DB                 *database.GormDatabase
	Provider           rp.RelyingParty
	UserChangeNotifier *UserChangeNotifier
	UsernameClaim      string
	PasswordStrength   int
	SecureCookie       bool
	AutoRegister       bool
	pendingSessions    *decaymap.DecayMap[string, *pendingOIDCSession]
}

// swagger:operation GET /auth/oidc/login oidc oidcLogin
//
// Start the OIDC login flow (browser).
//
// Redirects the user to the OIDC provider's authorization endpoint.
// After authentication, the provider redirects back to the callback endpoint.
//
//	---
//	parameters:
//	- name: name
//	  in: query
//	  description: the client name to create after login
//	  required: true
//	  type: string
//	responses:
//	  302:
//	    description: Redirect to OIDC provider
//	  default:
//	    description: Error
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *OIDCAPI) LoginHandler() gin.HandlerFunc {
	return gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		clientName := r.URL.Query().Get("name")
		if clientName == "" {
			http.Error(w, "invalid client name", http.StatusBadRequest)
			return
		}
		state, err := a.generateState()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to generate state: %v", err), http.StatusInternalServerError)
			return
		}
		a.pendingSessions.Set(time.Now(), state, &pendingOIDCSession{ClientName: clientName, CreatedAt: time.Now()})
		rp.AuthURLHandler(func() string { return state }, a.Provider)(w, r)
	})
}

// swagger:operation GET /auth/oidc/elevate oidc oidcElevate
//
// Start the OIDC flow to elevate an existing client session (browser).
//
// Redirects the user to the OIDC provider's authorization endpoint. After
// successful authentication, the referenced client session is elevated for
// the requested duration.
//
//	---
//	parameters:
//	- name: id
//	  in: query
//	  description: the client id to elevate
//	  required: true
//	  type: integer
//	  format: int64
//	- name: durationSeconds
//	  in: query
//	  description: how long the elevation should last, in seconds
//	  required: true
//	  type: integer
//	responses:
//	  302:
//	    description: Redirect to OIDC provider
//	  default:
//	    description: Error
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *OIDCAPI) ElevateHandler(ctx *gin.Context) {
	var elevate pendingElevation
	if err := ctx.BindQuery(&elevate); err != nil {
		return
	}
	state, err := a.generateState()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	a.pendingSessions.Set(time.Now(), state, &pendingOIDCSession{CreatedAt: time.Now(), Elevate: &elevate})
	rp.AuthURLHandler(func() string { return state }, a.Provider)(ctx.Writer, ctx.Request)
}

// swagger:operation GET /auth/oidc/callback oidc oidcCallback
//
// Handle the OIDC provider callback (browser).
//
// Exchanges the authorization code for tokens, resolves the user,
// creates a gotify client, sets a session cookie, and redirects to the UI.
//
//	---
//	parameters:
//	- name: code
//	  in: query
//	  description: the authorization code from the OIDC provider
//	  required: true
//	  type: string
//	- name: state
//	  in: query
//	  description: the state parameter for CSRF protection
//	  required: true
//	  type: string
//	responses:
//	  200:
//	    description: ok
//	  307:
//	    description: Redirect to UI
//	  default:
//	    description: Error
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *OIDCAPI) CallbackHandler() gin.HandlerFunc {
	callback := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, provider rp.RelyingParty, info *oidc.UserInfo) {
		user, status, err := a.resolveUser(info)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		session, ok := a.popPendingSession(state)
		if !ok {
			http.Error(w, "unknown or expired state", http.StatusBadRequest)
			return
		}

		if session.Elevate != nil {
			a.handleElevationCallback(w, session.Elevate, user)
			return
		}

		client, err := a.createClient(session.ClientName, user.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create client: %v", err), http.StatusInternalServerError)
			return
		}
		auth.SetCookie(w, client.Token, auth.CookieMaxAge, a.SecureCookie)
		// A reverse proxy may have already stripped a url prefix from the URL
		// without us knowing, we have to make a relative redirect.
		// We cannot use http.Redirect as this normalizes the Path with r.URL.
		w.Header().Set("Location", "../../")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
	return gin.WrapF(rp.CodeExchangeHandler(rp.UserinfoCallback(callback), a.Provider))
}

func (a *OIDCAPI) handleElevationCallback(w http.ResponseWriter, elevate *pendingElevation, user *model.User) {
	client, err := a.DB.GetClientByID(elevate.ClientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("database error: %v", err), http.StatusInternalServerError)
		return
	}
	if client == nil || client.UserID != user.ID {
		http.Error(w, "client not found", http.StatusNotFound)
		return
	}
	elevatedUntil := time.Now().Add(time.Duration(elevate.DurationSeconds) * time.Second)
	if err := a.DB.UpdateClientElevatedUntil(client.ID, &elevatedUntil); err != nil {
		http.Error(w, fmt.Sprintf("failed to elevate session: %v", err), http.StatusInternalServerError)
		return
	}

	// The UI rechecks the authentication when the tab is closed.
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, `<!DOCTYPE html>
<html lang="en">
<head>
  <title>Gotify Session Elevation</title>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
</head>
<body>
  <h1 style="text-align:center">Gotify session elevation successful. Close this tab to continue.</h1>
  <script>window.close();</script>
</body>
</html>`)
}

// swagger:operation POST /auth/oidc/external/authorize oidc externalAuthorize
//
// Initiate the OIDC authorization flow for a native app.
//
// The app generates a PKCE code_verifier and code_challenge, then calls this
// endpoint. The server forwards the code_challenge to the OIDC provider and
// returns the authorization URL for the app to open in a browser.
//
//	---
//	consumes: [application/json]
//	produces: [application/json]
//	parameters:
//	- name: body
//	  in: body
//	  required: true
//	  schema:
//	    $ref: "#/definitions/OIDCExternalAuthorizeRequest"
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/OIDCExternalAuthorizeResponse"
//	  default:
//	    description: Error
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *OIDCAPI) ExternalAuthorizeHandler(ctx *gin.Context) {
	var req model.OIDCExternalAuthorizeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	state, err := a.generateState()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	a.pendingSessions.Set(time.Now(), state, &pendingOIDCSession{
		RedirectURI: req.RedirectURI, ClientName: req.Name, CreatedAt: time.Now(),
	})
	authOpts := []rp.AuthURLOpt{
		rp.AuthURLOpt(rp.WithURLParam("redirect_uri", req.RedirectURI)),
		rp.WithCodeChallenge(req.CodeChallenge),
	}
	ctx.JSON(http.StatusOK, &model.OIDCExternalAuthorizeResponse{
		AuthorizeURL: rp.AuthURL(state, a.Provider, authOpts...),
		State:        state,
	})
}

// swagger:operation POST /auth/oidc/external/token oidc externalToken
//
// Exchange an authorization code for a gotify client token.
//
// After the user authenticates with the OIDC provider and the app receives
// the authorization code via redirect, the app calls this endpoint with the
// code and PKCE code_verifier. The server exchanges the code with the OIDC
// provider and returns a gotify client token.
//
//	---
//	consumes: [application/json]
//	produces: [application/json]
//	parameters:
//	- name: body
//	  in: body
//	  required: true
//	  schema:
//	    $ref: "#/definitions/OIDCExternalTokenRequest"
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/OIDCExternalTokenResponse"
//	  default:
//	    description: Error
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *OIDCAPI) ExternalTokenHandler(ctx *gin.Context) {
	var req model.OIDCExternalTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session, ok := a.popPendingSession(req.State)
	if !ok {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("unknown or expired state"))
		return
	}
	exchangeOpts := []rp.CodeExchangeOpt{
		rp.CodeExchangeOpt(rp.WithURLParam("redirect_uri", session.RedirectURI)),
		rp.WithCodeVerifier(req.CodeVerifier),
	}
	tokens, err := rp.CodeExchange[*oidc.IDTokenClaims](ctx.Request.Context(), req.Code, a.Provider, exchangeOpts...)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token exchange failed: %w", err))
		return
	}
	info, err := rp.Userinfo[*oidc.UserInfo](ctx.Request.Context(), tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.GetSubject(), a.Provider)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user info: %w", err))
		return
	}
	user, status, resolveErr := a.resolveUser(info)
	if resolveErr != nil {
		ctx.AbortWithError(status, resolveErr)
		return
	}
	client, err := a.createClient(session.ClientName, user.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, &model.OIDCExternalTokenResponse{
		Token: client.Token,
		User:  &model.UserExternal{ID: user.ID, Name: user.Name, Admin: user.Admin},
	})
}

func (a *OIDCAPI) generateState() (string, error) {
	nonce := make([]byte, 20)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil
}

// resolveUser looks up or creates a user from OIDC userinfo claims.
func (a *OIDCAPI) resolveUser(info *oidc.UserInfo) (*model.User, int, error) {
	usernameRaw, ok := info.Claims[a.UsernameClaim]
	if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("username claim %q is missing", a.UsernameClaim)
	}
	username := fmt.Sprint(usernameRaw)
	if username == "" || usernameRaw == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("username claim was empty")
	}

	user, err := a.DB.GetUserByName(username)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		if !a.AutoRegister {
			return nil, http.StatusForbidden, fmt.Errorf("user does not exist and auto-registration is disabled")
		}
		user = &model.User{Name: username, Admin: false, Pass: nil}
		if err := a.DB.CreateUser(user); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err)
		}
		if err := a.UserChangeNotifier.fireUserAdded(user.ID); err != nil {
			log.Printf("Could not notify user change: %v\n", err)
		}
	}
	return user, 0, nil
}

func (a *OIDCAPI) createClient(name string, userID uint) (*model.Client, error) {
	elevatedUntil := time.Now().Add(model.DefaultElevationDuration)
	client := &model.Client{
		Name:          name,
		Token:         auth.GenerateNotExistingToken(generateClientToken, func(t string) bool { c, _ := a.DB.GetClientByToken(t); return c != nil }),
		UserID:        userID,
		ElevatedUntil: &elevatedUntil,
	}
	return client, a.DB.CreateClient(client)
}

func (a *OIDCAPI) popPendingSession(key string) (*pendingOIDCSession, bool) {
	session, ok := a.pendingSessions.Pop(key)
	if ok && time.Since(session.CreatedAt) < pendingSessionMaxAge {
		return session, true
	}
	return nil, false
}
