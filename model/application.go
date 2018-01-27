package model

// Application holds information about an app which can send notifications.
type Application struct {
	ID          string    `gorm:"primary_key;unique_index"`
	UserID      uint      `gorm:"index" json:"-"`
	Name        string    `form:"name" query:"name" json:"name" binding:"required"`
	Description string    `form:"description" query:"description" json:"description"`
	Messages    []Message `json:"-"`
}
