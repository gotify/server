package database

import (
	"time"

	"github.com/gotify/server/v2/model"
	"gorm.io/gorm"
)

// GetClientByID returns the client for the given id or nil.
func (d *GormDatabase) GetClientByID(id uint) (*model.Client, error) {
	client := new(model.Client)
	err := d.notExpired(d.DB.Where("id = ?", id)).Find(client).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if client.ID == id {
		return client, err
	}
	return nil, err
}

// GetClientByToken returns the client for the given token or nil.
func (d *GormDatabase) GetClientByToken(token string) (*model.Client, error) {
	client := new(model.Client)
	err := d.notExpired(d.DB.Where("token = ?", token)).Find(client).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if client.Token == token {
		return client, err
	}
	return nil, err
}

// CreateClient creates a client.
func (d *GormDatabase) CreateClient(client *model.Client) error {
	if client.CreatedAt.IsZero() {
		client.CreatedAt = d.DB.NowFunc()
	}
	client.PopulateExpiresAt()
	return d.DB.Create(client).Error
}

// GetClientsByUser returns all clients from a user.
func (d *GormDatabase) GetClientsByUser(userID uint) ([]*model.Client, error) {
	var clients []*model.Client
	err := d.notExpired(d.DB.Where("user_id = ?", userID)).Find(&clients).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return clients, err
}

// DeleteClientByID deletes a client by its id.
func (d *GormDatabase) DeleteClientByID(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.Client{}).Error
}

// UpdateClient updates a client.
func (d *GormDatabase) UpdateClient(client *model.Client) error {
	client.PopulateExpiresAt()
	return d.DB.Save(client).Error
}

// UpdateClientTokensLastUsed updates the last used timestamp of clients and
// keeps the expires_at in sync.
func (d *GormDatabase) UpdateClientTokensLastUsedAndExpiresAt(tokens []string, t *time.Time) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		var clients []*model.Client
		if err := tx.Where("token IN (?)", tokens).Find(&clients).Error; err != nil {
			return err
		}
		for _, c := range clients {
			c.LastUsed = t
			c.PopulateExpiresAt()
			if err := tx.Save(c).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateClientElevatedUntil updates the elevated_until timestamp of a client by token.
func (d *GormDatabase) UpdateClientElevatedUntil(id uint, t *time.Time) error {
	return d.DB.Model(&model.Client{}).Where("id = ?", id).Update("elevated_until", t).Error
}

// CleanupExpiredClients deletes clients whose expires_at has passed.
func (d *GormDatabase) CleanupExpiredClients(now time.Time) ([]*model.Client, error) {
	var expired []*model.Client
	if err := d.DB.Where("expires_at IS NOT NULL AND expires_at <= ?", now).Find(&expired).Error; err != nil {
		return nil, err
	}
	if len(expired) == 0 {
		return nil, nil
	}
	ids := make([]uint, len(expired))
	for i, c := range expired {
		ids[i] = c.ID
	}
	if err := d.DB.Where("id IN ?", ids).Delete(&model.Client{}).Error; err != nil {
		return nil, err
	}
	return expired, nil
}

func (d *GormDatabase) notExpired(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", d.DB.NowFunc())
}
