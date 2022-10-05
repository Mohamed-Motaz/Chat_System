package Database

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	Common
	Token       string `gorm:"column:token"            json:"token"`
	Name        string `gorm:"column:name"             json:"name"`
	Chats_count int32  `gorm:"column:chats_count"      json:"chats_count"`
}

func (Application) TableName() string {
	return "instabug.applications"
}

func (db *DBWrapper) UpdateApplicationCtr(token string, chatsCount int, updatedAt time.Time) *gorm.DB {
	return db.Db.Exec(
		`
		UPDATE instabug.applications
		SET chats_count = ?, updated_at = ? 
		WHERE token = ?
		`, chatsCount, updatedAt, token)
}
