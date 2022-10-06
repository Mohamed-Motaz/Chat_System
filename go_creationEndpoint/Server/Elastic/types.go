package MessageQueue

import (
	utils "Server/Utils"
	"context"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
)

const (
	MESSAGES_INDEX string = "messagesIndex"
)

type Elastic struct {
	elastic *elasticsearch.Client
	ctx     context.Context
	timeout time.Duration
}

//elastic object
type ElasticObj struct {
	Id      int32  `json:"id"`
	Chat_id int32  `json:"chat_id"`
	Number  int32  `json:"number"`
	Body    string `json:"body"`
}

const (
	ELASTIC_PORT string = "ELASTIC_PORT"
	ELASTIC_HOST string = "ELASTIC_HOST"
	LOCAL_HOST   string = "127.0.0.1"
)

var ElasticHost string
var ElasticPort string

func init() {
	if !utils.IN_DOCKER {
		err := godotenv.Load("config.env")
		if err != nil {
			log.Fatal("Error loading config.env file")
		}
	}
	ElasticHost = strings.Replace(utils.GetEnv(ELASTIC_HOST, LOCAL_HOST), "_", ".", -1) //replace all "_" with ".""
	ElasticPort = utils.GetEnv(ELASTIC_PORT, "9200")

}
