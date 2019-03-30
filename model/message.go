package model

import (
	"time"
)

// Message holds information about a message
type Message struct {
	ID            uint `gorm:"AUTO_INCREMENT;primary_key;index"`
	ApplicationID uint
	Message       string `gorm:"type:text"`
	Title         string `gorm:"type:text"`
	Priority      int
	Extras        []byte
	Date          time.Time
}

// MessageExternal Model
//
// The MessageExternal holds information about a message which was sent by an Application.
//
// swagger:model Message
type MessageExternal struct {
	// The message id.
	//
	// read only: true
	// required: true
	// example: 25
	ID uint `json:"id"`
	// The application id that send this message.
	//
	// read only: true
	// required: true
	// example: 5
	ApplicationID uint `json:"appid"`
	// The message. Markdown (excluding html) is allowed.
	//
	// required: true
	// example: **Backup** was successfully finished.
	Message string `form:"message" query:"message" json:"message" binding:"required"`
	// The title of the message.
	//
	// example: Backup
	Title string `form:"title" query:"title" json:"title"`
	// The priority of the message.
	//
	// example: 2
	Priority int `form:"priority" query:"priority" json:"priority"`
	// The extra data sent along the message.
	//
	// The extra fields are stored in a key-value scheme. Only accepted in CreateMessage requests with application/json content-type.
	//
	// The keys should be in the following format: &lt;top-namespace&gt;::[&lt;sub-namespace&gt;::]&lt;action&gt;
	//
	// These namespaces are reserved and might be used in the official clients: gotify android ios web server client. Do not use them for other purposes.
	//
	// example: {"home::appliances::thermostat::change_temperature":{"temperature":23},"home::appliances::lighting::on":{"brightness":15}}
	Extras map[string]interface{} `form:"-" query:"-" json:"extras,omitempty"`
	// The date the message was created.
	//
	// read only: true
	// required: true
	// example: 2018-02-27T19:36:10.5045044+01:00
	Date time.Time `json:"date"`
}
