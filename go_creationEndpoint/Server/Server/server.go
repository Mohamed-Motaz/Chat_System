package Server

import (
	db "Server/Database"
	logger "Server/Logger"

	"net/http"

	"github.com/rs/cors"
)

//create a new server
func NewServer() (*Server, error) {

	server := &Server{
		dBWrapper: db.New(),
	}

	mux := http.NewServeMux()
	registerRoutes(mux)
	server.handler = cors.AllowAll().Handler(mux)

	go server.serve()

	return server, nil
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/getAllProducts", nil)

}

func (server *Server) serve() {
	logger.LogInfo(logger.SERVER, logger.ESSENTIAL, "Listening on %v:%v", MyHost, MyPort)
	err := http.ListenAndServe(MyHost+":"+MyPort, server.handler)

	if err != nil {
		logger.FailOnError(logger.SERVER, logger.ESSENTIAL, "Failed in listening on port with error %v", err)
	}
}
