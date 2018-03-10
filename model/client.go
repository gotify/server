package model

// Client Model
//
// The Client holds information about a device which can receive notifications (and other stuff).
//
// swagger:model Client
type Client struct {
	// The client id.
	//
	// read only: true
	// required: true
	// example: 5
	ID uint `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"id"`
	// The client token. Can be used as `clientToken`. See Authentication.
	//
	// read only: true
	// required: true
	// example: CWH0wZ5r0Mbac.r
	Token  string `gorm:"unique_index" json:"token"`
	UserID uint   `gorm:"index" json:"-"`
	// The client name. This is how the client should be displayed to the user.
	//
	// required: true
	// example: Android Phone
	Name string `form:"name" query:"name" json:"name" binding:"required"`
}
