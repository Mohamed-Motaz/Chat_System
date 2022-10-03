package Database

import (
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

func (db *DBWrapper) GetApplicationByToken(a *Application, token string) *gorm.DB {
	return db.Db.Raw("SELECT * FROM instabug.applications WHERE token = ?", token).Scan(a)
}

func (db *DBWrapper) UpdateApplicationByToken(a *Application, name, token string) *gorm.DB {
	return db.Db.Raw(`UPDATE instabug.applications 
					  SET name = ? 
					WHERE token = ?`, name, token).Scan(a)
}

func (db *DBWrapper) InsertApplication(a *Application) *gorm.DB {
	return db.Db.Create(a)
}

func (db *DBWrapper) UpdateApplication(a *Application) *gorm.DB {
	return db.Db.Save(a)
}
