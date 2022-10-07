package Worker

import (
	db "Worker/Database"
	es "Worker/Elastic"
	q "Worker/MessageQueue"

	utils "Worker/Utils"

	"log"
	"strings"

	"github.com/joho/godotenv"
)

type Worker struct {
	dBWrapper *db.DBWrapper
	Mq        *q.MQ
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
	MyPort = utils.GetEnv(MY_PORT, "6666")

}
