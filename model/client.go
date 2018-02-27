package model

// Client Model
//
// The Client holds information about a device which can receive notifications (and other stuff).
//
// swagger:model Client
type Client struct {
	// The client id. Can be used as `clientToken`. See Authentication.
	//
	// read only: true
	// required: true
	// example: CWH0wZ5r0Mbac.r
	ID     string `gorm:"primary_key;unique_index" json:"id"`
	UserID uint   `gorm:"index" json:"-"`
	// The client name. This is how the client should be displayed to the user.
	//
	// required: true
	// example: Android Phone
	Name   string `form:"name" query:"name" json:"name" binding:"required"`
}
