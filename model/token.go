package model

type Token struct {
	Id           string `gorm:"primary_key;unique_index"`
	UserID       uint   `gorm:"index" json:"-"`
	Name         string `form:"name" query:"name" json:"name" binding:"required"`
	Description  string `form:"description" query:"description" json:"description"`
	WriteOnly    bool `form:"writeOnly" query:"writeOnly" json:"writeOnly" binding:"exists"`
	Messages     []Message `json:"-"`
}
