package model

// Application Model
//
// The Application holds information about an app which can send notifications.
//
// swagger:model Application
type Application struct {
	// The application id. Can be used as `appToken`. See Authentication.
	//
	// read only: true
	// required: true
	// example: AWH0wZ5r0Mbac.r
	ID          string    `gorm:"primary_key;unique_index" json:"id"`
	UserID      uint      `gorm:"index" json:"-"`
	// The application name. This is how the application should be displayed to the user.
	//
	// required: true
	// example: Backup Server
	Name        string    `form:"name" query:"name" json:"name" binding:"required"`
	// The description of the application.
	//
	// required: true
	// example: Backup server for the interwebs
	Description string    `form:"description" query:"description" json:"description"`
	Messages    []Message `json:"-"`
}
