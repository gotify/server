package model

// Application Model
//
// The Application holds information about an app which can send notifications.
//
// swagger:model Application
type Application struct {
	// The application id.
	//
	// read only: true
	// required: true
	// example: 5
	ID uint `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"id"`
	// The application token. Can be used as `appToken`. See Authentication.
	//
	// read only: true
	// required: true
	// example: AWH0wZ5r0Mbac.r
	Token  string `gorm:"type:varchar(180);unique_index" json:"token"`
	UserID uint   `gorm:"index" json:"-"`
	// The application name. This is how the application should be displayed to the user.
	//
	// required: true
	// example: Backup Server
	Name string `gorm:"type:text" form:"name" query:"name" json:"name" binding:"required"`
	// The description of the application.
	//
	// required: true
	// example: Backup server for the interwebs
	Description string `gorm:"type:text" form:"description" query:"description" json:"description"`
	// Whether the application is an internal application. Internal applications should not be deleted.
	//
	// read only: true
	// required: true
	// example: false
	Internal bool `form:"internal" query:"internal" json:"internal"`
	// The image of the application.
	//
	// read only: true
	// required: true
	// example: image/image.jpeg
	Image    string            `gorm:"type:text" json:"image"`
	Messages []MessageExternal `json:"-"`
	// The priority default of a message for an application. If the sender does not provide a priority, then
	// the message priority will be set to this value. Defaults to 0.
	//
	// required: false
	// example: 4
	PriorityDefault int `gorm:"default:0" form:"priorityDefault" query:"priorityDefault" json:"priorityDefault"`
}
