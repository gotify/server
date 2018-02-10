package model

// The User holds information about the credentials of a user and its application and client tokens.
type User struct {
	ID           uint `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name         string `gorm:"unique_index"`
	Pass         []byte
	Admin        bool
	Applications []Application
	Clients      []Client
}
