package model

type Message struct {
	ID       uint `gorm:"primary_key" gorm:"AUTO_INCREMENT;primary_key;index"`
	TokenID  string
	Message  string
	Title    string
	Priority int
}
