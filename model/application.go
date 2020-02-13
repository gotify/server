package model

import (
	"encoding/json"
	"time"
)

// TimeNow is used to retrieve the current time. This variable was introduced to ease testing.
var TimeNow = time.Now

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
	Image    string               `gorm:"type:text" json:"image"`
	Messages []ApplicationMessage `json:"-"`
}

// ApplicationMessage Model
//
// The ApplicationMessage holds information about a message which was sent by an application.
//
// swagger:model ApplicationMessage
type ApplicationMessage struct {
	// The message. Markdown (excluding HTML) is allowed.
	//
	// required: true
	// example: **Backup** was successfully finished.
	Message string `form:"message" json:"message" binding:"required"`
	// The title of the message.
	//
	// example: Backup
	Title string `form:"title" json:"title"`
	// The priority of the message.
	//
	// example: 2
	Priority int `form:"priority" json:"priority"`
	// The extra data sent along the message.
	//
	// The extra fields are stored in a key-value scheme. Only accepted in CreateMessage requests with application/json content-type.
	//
	// The keys should be in the following format: &lt;top-namespace&gt;::[&lt;sub-namespace&gt;::]&lt;action&gt;
	//
	// These namespaces are reserved and might be used in the official clients: gotify android ios web server client. Do not use them for other purposes.
	//
	// example: {"home::appliances::thermostat::change_temperature":{"temperature":23},"home::appliances::lighting::on":{"brightness":15}}
	Extras map[string]interface{} `form:"-" json:"extras,omitempty"`
}

// ToInternal converts the ApplicationMessage to an internal representation.
func (msg ApplicationMessage) ToInternal(applicationID uint) *Message {
	res := &Message{
		Message:       msg.Message,
		Title:         msg.Title,
		Priority:      msg.Priority,
	}
	res.ApplicationID = applicationID
	res.Date = TimeNow()
	if msg.Extras != nil {
		res.Extras, _ = json.Marshal(msg.Extras)
	}
	return res
}
