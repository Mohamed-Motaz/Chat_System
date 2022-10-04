package Database

type Message struct {
	Common
	Chat_id int32  `gorm:"column:chat_id"       json:"chat_id"`
	Number  int32  `gorm:"column:number"        json:"number"`
	Body    string `gorm:"column:body"          json:"body"`
}

func (Message) TableName() string {
	return "messages"
}
