package model

// SecurityUpdateAction Model
//
// The SecurityUpdateAction describes the details of a requested security update.
//
// swagger:model SecurityUpdateAction
type SecurityUpdateAction struct {
	// Whether to regenerate the token. Your client token must be elevated to perform this action.
	//
	// example: true
	RegenerateToken bool `form:"regenerateToken" query:"regenerateToken" json:"regenerateToken"`
}

// SecurityUpdateActionResponse Model
//
// The SecurityUpdateActionResponse holds information about the response to a security update request.
//
// swagger:model SecurityUpdateActionResponse
type SecurityUpdateActionResponse struct {
	// The response to the regenerate token action. Only present if the regenerate token action was requested.
	RegenerateToken *RegenerateTokenResponse `json:"regenerateToken,omitempty"`
}

// RegenerateTokenResponse Model
//
// The RegenerateTokenResponse holds information about the response to the regenerate token action.
//
// swagger:model RegenerateTokenResponse
type RegenerateTokenResponse struct {
	// The new token.
	//
	// example: gtfya.e2NcJK7AenXBPIRB3S03JsBlmy0V6xP8h0hwSiAJae8
	// read only: true
	// required: true
	Token string `json:"token"`
}
