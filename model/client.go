package model

type Client struct {
	Id     string `gorm:"primary_key;unique_index"`
	UserId uint   `gorm:"index" json:"-"`
	Name   string `form:"name" query:"name" json:"name" binding:"required"`
}
