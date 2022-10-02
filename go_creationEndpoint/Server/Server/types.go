package Server

import (
	db "Server/Database"
	utils "Server/Utils"

	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

type Server struct {
	dBWrapper *db.DBWrapper
	handler   http.Handler //custom handler that logs requests
}

const (
	MY_PORT    string = "MY_PORT"
	MY_HOST    string = "MY_HOST"
	LOCAL_HOST string = "127.0.0.1"
)

var MyHost string
var MyPort string
var DebugEnv int16

func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env file")
	}
	MyHost = strings.Replace(utils.GetEnv(MY_HOST, LOCAL_HOST), "_", ".", -1) //replace all "_" with ".""
	MyPort = utils.GetEnv(MY_PORT, "5555")

}
