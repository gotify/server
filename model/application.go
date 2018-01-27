package model

type Application struct {
	Id          string    `gorm:"primary_key;unique_index"`
	UserId      uint      `gorm:"index" json:"-"`
	Name        string    `form:"name" query:"name" json:"name" binding:"required"`
	Description string    `form:"description" query:"description" json:"description"`
	Messages    []Message `json:"-"`
}
