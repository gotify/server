package model

import (
	"encoding/json"
	"time"
)

// Message holds information about a new message.
type Message struct {
	ID            uint `gorm:"AUTO_INCREMENT;primary_key;index"`
	ApplicationID uint
	Message       string `gorm:"type:text"`
	Title         string `gorm:"type:text"`
	Priority      int
	Extras        []byte
	Date          time.Time
}

// ToExternal converts the event into an external representation.
func (msg Message) ToExternal() interface{} {
	res := &MessageExternal{
		ID:            msg.ID,
		ApplicationID: msg.ApplicationID,
		Message:       msg.Message,
		Title:         msg.Title,
		Priority:      msg.Priority,
		Date:          msg.Date,
	}
	if len(msg.Extras) != 0 {
		res.Extras = make(map[string]interface{})
		json.Unmarshal(msg.Extras, &res.Extras)
	}
	return res
}

// MessageExternal Model
//
// MessageExternal holds information about a message which will be sent to the clients.
//
// swagger:model Message
type MessageExternal struct {
	// The message id.
	//
	// read only: true
	// required: true
	// example: 25
	ID uint `json:"id"`
	// The ID of the application that sent this message.
	//
	// read only: true
	// required: true
	// example: 5
	ApplicationID uint `json:"appid"`
	// The message. Markdown (excluding HTML) is allowed.
	//
	// required: true
	// example: **Backup** was successfully finished.
	Message string `json:"message"`
	// The title of the message.
	//
	// example: Backup
	Title string `json:"title"`
	// The priority of the message.
	//
	// example: 2
	Priority int `json:"priority"`
	// The extra data sent along the message.
	//
	// The extra fields are stored in a key-value scheme. Only accepted in CreateMessage requests with application/json content-type.
	//
	// The keys should be in the following format: &lt;top-namespace&gt;::[&lt;sub-namespace&gt;::]&lt;action&gt;
	//
	// These namespaces are reserved and might be used in the official clients: gotify android ios web server client. Do not use them for other purposes.
	//
	// example: {"home::appliances::thermostat::change_temperature":{"temperature":23},"home::appliances::lighting::on":{"brightness":15}}
	Extras map[string]interface{} `json:"extras,omitempty"`
	// The date the message was created.
	//
	// read only: true
	// required: true
	// example: 2018-02-27T19:36:10.5045044+01:00
	Date time.Time `json:"date"`
}
