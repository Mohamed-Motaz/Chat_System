package Database

import "gorm.io/gorm"

type Chat struct {
	Common
	Application_token string `gorm:"column:application_token"  json:"application_token"`
	Number            int32  `gorm:"column:number"             json:"number"`
	Messages_count    int32  `gorm:"column:messages_count"     json:"messages_count"`
}

func (Chat) TableName() string {
	return "chats"
}

func (db *DBWrapper) GetChatByAppTokenAndNumber(c *Chat, token string, chatNum int) *gorm.DB {
	return db.Db.Raw(`SELECT * FROM instabug.chats 
					  WHERE application_token = ?
					  AND number = ? LIMIT 1`, token, chatNum).Scan(c)
}

func (db *DBWrapper) InsertChat(c *Chat) *gorm.DB {
	return db.Db.Create(c)
}

func (db *DBWrapper) UpdateChat(c *Chat) *gorm.DB {
	return db.Db.Save(c)
}
