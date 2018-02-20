package model

// VersionInfo Model
//
// swagger:model VersionInfo
type VersionInfo struct {
	Version string `json:"version"`
	Commit string `json:"commit"`
	BuildDate string `json:"buildDate"`
	Branch string `json:"branch"`
}

