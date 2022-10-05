package Database

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	Common
	Application_token string `gorm:"column:application_token"  json:"application_token"`
	Number            int32  `gorm:"column:number"             json:"number"`
	Messages_count    int32  `gorm:"column:messages_count"     json:"messages_count"`
}

func (Chat) TableName() string {
	return "instabug.chats"
}

func (db *DBWrapper) UpdateChatsCtr(token string, number int, messagesCount int, updatedAt time.Time) *gorm.DB {
	return db.Db.Exec(
		`
		UPDATE instabug.chats
		SET messages_count = ?, updated_at = ? 
		WHERE application_token = ? and number = ?
		`, messagesCount, updatedAt, token, number)
}

func (db *DBWrapper) InsertChat(c *Chat) *gorm.DB {
	return db.Db.Create(c)
}

func (db *DBWrapper) UpdateChat(c *Chat) *gorm.DB {
	return db.Db.Save(c)
}
