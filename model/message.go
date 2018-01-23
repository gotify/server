package model

type Message struct {
	ID       uint `gorm:"AUTO_INCREMENT;primary_key;index"`
	TokenID  string
	Message  string
	Title    string
	Priority int
}
