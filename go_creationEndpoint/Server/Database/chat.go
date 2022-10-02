package Database

type Chat struct {
	Common
	Application_id int32 `gorm:"column:application_id"        	   json:"application_id"`
	Number         int32 `gorm:"column:number"             json:"number"`
	Messages_count int32 `gorm:"column:messages_count"      json:"messages_count"`
}

func (Chat) TableName() string {
	return "chats"
}
