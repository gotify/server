package model

// GotifyInfo Model
//
// swagger:model GotifyInfo
type GotifyInfo struct {
	// The current version.
	//
	// required: true
	// example: 5.2.6
	Version string `json:"version"`
	// If registration is enabled.
	//
	// required: true
	// example: true
	Register bool `json:"register"`
	// If local authentication is enabled.
	//
	// required: true
	// example: true
	LocalAuth bool `json:"localAuth"`
	// If oidc is enabled.
	//
	// required: true
	// example: true
	Oidc bool `json:"oidc"`
	// Name of the OIDC identity provider.
	//
	// required: true
	// example: OIDC
	OIDCIDPName string `json:"oidcIdpName"`
	// If the WebUI should automatically redirect to the OIDC identity
	// provider instead of showing the login page.
	//
	// required: true
	// example: false
	OIDCAutoRedirect bool `json:"oidcAutoRedirect"`
}
