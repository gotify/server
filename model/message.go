package model

import "time"

// The Message holds information about a message which was sent by an Application.
type Message struct {
	ID            uint      `gorm:"AUTO_INCREMENT;primary_key;index" json:"id"`
	ApplicationID string    `json:"appid"`
	Message       string    `form:"message" query:"message" json:"message" binding:"required"`
	Title         string    `form:"title" query:"title" json:"title" binding:"required"`
	Priority      int       `form:"priority" query:"priority" json:"priority"`
	Date          time.Time `json:"date"`
}
