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
	// If oidc is enabled.
	//
	// required: true
	// example: true
	Oidc bool `json:"oidc"`
}
