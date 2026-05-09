package model

import "time"

// Client Model
//
// The Client holds information about a device which can receive notifications (and other stuff).
//
// swagger:model Client
type Client struct {
	// The client id.
	//
	// read only: true
	// required: true
	// example: 5
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// The client token. Can be used as `clientToken`. See Authentication.
	//
	// read only: true
	// required: true
	// example: CWH0wZ5r0Mbac.r
	Token  string `gorm:"type:varchar(180);uniqueIndex:uix_clients_token" json:"token"`
	UserID uint   `gorm:"index" json:"-"`
	// The client name. This is how the client should be displayed to the user.
	//
	// required: true
	// example: Android Phone
	Name string `gorm:"type:text" form:"name" query:"name" json:"name" binding:"required"`
	// The date the client was created.
	//
	// read only: true
	// required: true
	// example: 2019-01-01T00:00:00Z
	CreatedAt time.Time `json:"createdAt"`
	// The last time the client token was used.
	//
	// read only: true
	// example: 2019-01-01T00:00:00Z
	LastUsed *time.Time `json:"lastUsed"`
	// The time until which this client's session is elevated.
	//
	// read only: true
	ElevatedUntil *time.Time `json:"elevatedUntil,omitempty"`
	// The number of seconds of inactivity after which the client is removed.
	// 0 means the client never expires.
	//
	// example: 2592000
	ExpiresAfterInactivitySeconds uint `gorm:"default:0;not null" form:"expiresAfterInactivitySeconds" query:"expiresAfterInactivitySeconds" json:"expiresAfterInactivitySeconds"`
	// The time at which this client will expire due to inactivity, or null if it never expires.
	//
	// read only: true
	// example: 2019-01-01T00:00:00Z
	ExpiresAt *time.Time `gorm:"index" json:"expiresAt,omitempty"`
}

func (c *Client) PopulateExpiresAt() {
	c.ExpiresAt = c.calculateExpiresAt()
}

func (c *Client) calculateExpiresAt() *time.Time {
	if c.ExpiresAfterInactivitySeconds == 0 {
		return nil
	}
	reference := c.CreatedAt
	if c.LastUsed != nil {
		reference = *c.LastUsed
	}
	expiry := reference.Add(time.Duration(c.ExpiresAfterInactivitySeconds) * time.Second)
	return &expiry
}
