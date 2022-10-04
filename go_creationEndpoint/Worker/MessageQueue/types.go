package MessageQueue

import (
	db "Worker/Database"
	utils "Worker/Utils"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

//queue names
const ENTITIES_QUEUE = "entitiesQueue"

type MQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	qMap map[string]*amqp.Queue
	mu   sync.Mutex
}

//objects passed into the messageQ

type Chat struct {
	db.Chat
	Id         int       `gorm:"primaryKey; column:id"   json:"id"`
	Created_at time.Time `gorm:"column:created_at"       json:"created_at"`
	Updated_at time.Time `gorm:"column:updated_at"       json:"updated_at"`
}

type Message struct {
	db.Message
	Id         int       `gorm:"primaryKey; column:id"   json:"id"`
	Created_at time.Time `gorm:"column:created_at"       json:"created_at"`
	Updated_at time.Time `gorm:"column:updated_at"       json:"updated_at"`
}

const (
	MQ_PORT     string = "MQ_PORT"
	MQ_HOST     string = "MQ_HOST"
	LOCAL_HOST  string = "127.0.0.1"
	MQ_USERNAME string = "MQ_USERNAME"
	MQ_PASSWORD string = "MQ_PASSWORD"
)

var MqHost string
var MqPort string
var MqUsername string
var MqPassword string

func init() {
	if !utils.IN_DOCKER {
		err := godotenv.Load("config.env")
		if err != nil {
			log.Fatal("Error loading config.env file")
		}
	}
	MqHost = strings.Replace(utils.GetEnv(MQ_HOST, LOCAL_HOST), "_", ".", -1) //replace all "_" with ".""
	MqPort = utils.GetEnv(MQ_PORT, "5672")
	MqUsername = utils.GetEnv(MQ_USERNAME, "guest")
	MqPassword = utils.GetEnv(MQ_PASSWORD, "guest")

}
