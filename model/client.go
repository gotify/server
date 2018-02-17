package model

// Client Model
//
// The Client holds information about a device which can receive notifications (and other stuff).
//
// swagger:model Client
type Client struct {
	ID     string `gorm:"primary_key;unique_index" json:"id"`
	UserID uint   `gorm:"index" json:"-"`
	Name   string `form:"name" query:"name" json:"name" binding:"required"`
}
