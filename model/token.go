package model

type Token struct {
	Name         string
	DefaultTitle string
	Description  string
	Icon         string
	WriteOnly    bool
	UserID       uint   `gorm:"index"`
	Id           string `gorm:"primary_key;unique_index"`
	Messages     []Message
}
