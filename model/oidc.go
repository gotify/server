package model

// OIDCExternalAuthorizeRequest Model
//
// Used to initiate the OIDC authorization flow for an external client.
//
// swagger:model OIDCExternalAuthorizeRequest
type OIDCExternalAuthorizeRequest struct {
	// The PKCE code challenge (S256).
	//
	// required: true
	// example: E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM
	CodeChallenge string `json:"code_challenge" binding:"required"`
	// The app's redirect URI.
	//
	// required: true
	// example: gotify://oidc/callback
	RedirectURI string `json:"redirect_uri" binding:"required"`
	// The client name to display in gotify.
	//
	// required: true
	// example: Android Phone
	Name string `json:"name" binding:"required"`
}

// OIDCExternalAuthorizeResponse Model
//
// Returned after initiating the OIDC authorization flow.
//
// swagger:model OIDCExternalAuthorizeResponse
type OIDCExternalAuthorizeResponse struct {
	// The URL to open in the browser to authenticate with the OIDC provider.
	//
	// required: true
	// example: https://auth.example.com/authorize?client_id=gotify&...
	AuthorizeURL string `json:"authorize_url"`
	// The state parameter to send back with the token exchange request.
	//
	// required: true
	// example: Android Phone:a1b2c3d4e5f6
	State string `json:"state"`
}

// OIDCExternalTokenRequest Model
//
// Used to exchange an authorization code for a gotify client token.
//
// swagger:model OIDCExternalTokenRequest
type OIDCExternalTokenRequest struct {
	// The authorization code from the OIDC provider.
	//
	// required: true
	Code string `json:"code" binding:"required"`
	// The state from the authorize response.
	//
	// required: true
	// example: Android Phone:a1b2c3d4e5f6
	State string `json:"state" binding:"required"`
	// The PKCE code verifier.
	//
	// required: true
	// example: dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk
	CodeVerifier string `json:"code_verifier" binding:"required"`
}

// OIDCExternalTokenResponse Model
//
// Returned after a successful token exchange.
//
// swagger:model OIDCExternalTokenResponse
type OIDCExternalTokenResponse struct {
	// The gotify client token for API authentication.
	//
	// required: true
	// example: CWH0wZ5r0Mbac.r
	Token string `json:"token"`
	// The authenticated user.
	//
	// required: true
	User *UserExternal `json:"user"`
}
