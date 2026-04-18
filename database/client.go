package database

import (
	"time"

	"github.com/gotify/server/v2/model"
	"gorm.io/gorm"
)

// GetClientByID returns the client for the given id or nil.
func (d *GormDatabase) GetClientByID(id uint) (*model.Client, error) {
	client := new(model.Client)
	err := d.DB.Where("id = ?", id).Find(client).Error
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
	err := d.DB.Where("token = ?", token).Find(client).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if client.Token == token {
		return client, err
	}
	return nil, err
}

func (d *GormDatabase) CountClientsByUserID(tx *gorm.DB, userID uint) (int64, error) {
	var count int64
	err := tx.Model(&model.Client{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CreateClient creates a client.
func (d *GormDatabase) CreateClient(client *model.Client, quota uint32) error {
	txn := d.DB.Begin()
	defer txn.Rollback()
	res := txn.Create(client)
	if res.Error != nil {
		return res.Error
	}
	if quota > 0 {
		count, err := d.CountClientsByUserID(txn, client.UserID)
		if err != nil {
			return err
		}
		if uint64(count) > uint64(quota) {
			// quota exceeded, delete the oldest client
			var oldestClient model.Client
			err := txn.Where("user_id = ?", client.UserID).Order("last_used ASC").First(&oldestClient).Error
			qe := ErrQuotaExceeded
			if err != nil {
				return qe
			}
			err = txn.Delete(&oldestClient).Error
			if err != nil {
				return qe
			}
		}
	}
	return txn.Commit().Error
}

// GetClientsByUser returns all clients from a user.
func (d *GormDatabase) GetClientsByUser(userID uint) ([]*model.Client, error) {
	var clients []*model.Client
	err := d.DB.Where("user_id = ?", userID).Find(&clients).Error
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
	return d.DB.Save(client).Error
}

// UpdateClientTokensLastUsed updates the last used timestamp of clients.
func (d *GormDatabase) UpdateClientTokensLastUsed(tokens []string, t *time.Time) error {
	return d.DB.Model(&model.Client{}).Where("token IN (?)", tokens).Update("last_used", t).Error
}
