package Server

import (
	db "Server/Database"
	es "Server/Elastic"
	q "Server/MessageQueue"
	c "Server/RedisCache"

	utils "Server/Utils"

	"log"
	"strings"

	"github.com/joho/godotenv"
)

type Server struct {
	dBWrapper *db.DBWrapper
	Mq        *q.MQ
	cache     *c.Cache
	elastic   *es.Elastic
}

const (
	MY_PORT    string = "MY_PORT"
	MY_HOST    string = "MY_HOST"
	LOCAL_HOST string = "127.0.0.1"
)

var MyHost string
var MyPort string

func init() {
	if !utils.IN_DOCKER {
		err := godotenv.Load("config.env")
		if err != nil {
			log.Fatal("Error loading config.env file")
		}
	}
	MyHost = strings.Replace(utils.GetEnv(MY_HOST, LOCAL_HOST), "_", ".", -1) //replace all "_" with ".""
	MyPort = utils.GetEnv(MY_PORT, "5555")

}
