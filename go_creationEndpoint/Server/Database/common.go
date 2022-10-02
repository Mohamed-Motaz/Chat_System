package Database

import "time"

type Common struct {
	Id         int       `gorm:"primaryKey; column:id"   json:"-"`
	Created_at time.Time `gorm:"column:created_at"       json:"created_at"`
	Updated_at time.Time `gorm:"column:updated_at"       json:"updated_at"`
}
