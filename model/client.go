package model

// The Client holds information about a device which can receive notifications (and other stuff).
type Client struct {
	ID     string `gorm:"primary_key;unique_index"`
	UserID uint   `gorm:"index" json:"-"`
	Name   string `form:"name" query:"name" json:"name" binding:"required"`
}
