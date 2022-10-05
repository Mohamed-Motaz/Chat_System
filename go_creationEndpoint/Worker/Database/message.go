package Database

import "gorm.io/gorm"

type Message struct {
	Common
	Chat_id int32  `gorm:"column:chat_id"       json:"chat_id"`
	Number  int32  `gorm:"column:number"        json:"number"`
	Body    string `gorm:"column:body"          json:"body"`
}

func (Message) TableName() string {
	return "instabug.messages"
}

func (db *DBWrapper) InsertMessage(m *Message) *gorm.DB {
	return db.Db.Create(m)
}

func (db *DBWrapper) UpdateMessage(m *Message) *gorm.DB {
	return db.Db.Save(m)
}
