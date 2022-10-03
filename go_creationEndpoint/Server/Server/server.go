package Server

import (
	db "Server/Database"
	logger "Server/Logger"
	q "Server/MessageQueue"

	"net/http"

	"github.com/gorilla/mux"
)

//create a new server
func NewServer() (*Server, error) {

	server := &Server{
		dBWrapper: db.New(),
		Mq:        *q.New("amqp://" + q.MqUsername + ":" + q.MqPassword + "@" + q.MqHost + ":" + q.MqPort + "/"),
	}
	r := mux.NewRouter()
	registerRoutes(r, server)

	go server.serve(r)

	return server, nil
}

func registerRoutes(r *mux.Router, s *Server) {
	r.HandleFunc("/api/v1/applications/{appToken}/chats", s.addChat).Methods("POST")

}

func (server *Server) serve(r *mux.Router) {
	logger.LogInfo(logger.SERVER, logger.ESSENTIAL, "Listening on %v:%v", MyHost, MyPort)
	err := http.ListenAndServe(MyHost+":"+MyPort, r)

	if err != nil {
		logger.FailOnError(logger.SERVER, logger.ESSENTIAL, "Failed in listening on port with error %v", err)
	}
}
