package model

type User struct {
	ID     uint `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name   string
	Pass   []byte
	Admin  bool
	Tokens []Token
}
