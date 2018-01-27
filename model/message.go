package model

import "time"

// The Message holds information about a message which was sent by an Application.
type Message struct {
	ID       uint `gorm:"AUTO_INCREMENT;primary_key;index"`
	TokenID  string
	Message  string
	Title    string
	Priority int
	Date     time.Time
}
