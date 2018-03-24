package model

// VersionInfo Model
//
// swagger:model VersionInfo
type VersionInfo struct {
	// The current version.
	//
	// required: true
	// example: 5.2.6
	Version string `json:"version"`
	// The git commit hash on which this binary was built.
	//
	// required: true
	// example: ae9512b6b6feea56a110d59a3353ea3b9c293864
	Commit string `json:"commit"`
	// The date on which this binary was built.
	//
	// required: true
	// example: 2018-02-27T19:36:10.5045044+01:00
	BuildDate string `json:"buildDate"`
}
