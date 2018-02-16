package model

// The User holds information about the credentials of a user and its application and client tokens.
type User struct {
	ID           uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name         string `gorm:"unique_index"`
	Pass         []byte
	Admin        bool
	Applications []Application
	Clients      []Client
}

// UserExternal Model
//
// The User holds information about the credentials and other stuff.
//
// swagger:model User
type UserExternal struct {
	ID    uint   `json:"id"`
	Name  string `binding:"required" json:"name" query:"name" form:"name"`
	Pass  string `json:"pass,omitempty" form:"pass" query:"pass"`
	Admin bool   `json:"admin" form:"admin" query:"admin"`
}
