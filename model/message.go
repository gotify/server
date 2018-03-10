package model

import "time"

// Message Model
//
// The Message holds information about a message which was sent by an Application.
//
// swagger:model Message
type Message struct {
	// The message id.
	//
	// read only: true
	// required: true
	// example: 25
	ID uint `gorm:"AUTO_INCREMENT;primary_key;index" json:"id"`
	// The application id that send this message.
	//
	// read only: true
	// required: true
	// example: 5
	ApplicationID uint `json:"appid"`
	// The actual message.
	//
	// required: true
	// example: Backup was successfully finished.
	Message string `form:"message" query:"message" json:"message" binding:"required"`
	// The title of the message.
	//
	// required: true
	// example: Backup
	Title string `form:"title" query:"title" json:"title" binding:"required"`
	// The priority of the message.
	//
	// example: 2
	Priority int `form:"priority" query:"priority" json:"priority"`
	// The date the message was created.
	//
	// read only: true
	// required: true
	// example: 2018-02-27T19:36:10.5045044+01:00
	Date time.Time `json:"date"`
}
