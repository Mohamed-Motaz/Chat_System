package Database

type Application struct {
	Common
	Token       string `gorm:"column:token"            json:"token"`
	Name        string `gorm:"column:name"             json:"name"`
	Chats_count int32  `gorm:"column:chats_count"      json:"chats_count"`
}

func (Application) TableName() string {
	return "applications"
}
