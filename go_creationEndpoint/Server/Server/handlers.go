package Server

import (
	db "Server/Database"
	logger "Server/Logger"
	q "Server/MessageQueue"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (server *Server) addApplication(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		failure(w, r, http.StatusBadRequest, "invalid request")
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "unable to read request")
		return
	}

	req := &AddApplicationReq{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to parse request to addApplication %v\nwith error %v", string(bytes), err)
		failure(w, r, http.StatusBadRequest, "unable to martial request")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "No name in request to addApplication %v", err)
		failure(w, r, http.StatusBadRequest, "no valid name")
		return
	}

	logger.LogInfo(logger.SERVER, logger.NON_ESSENTIAL, "The http received message to addApplication is %+v", req)

	a := &db.Application{
		Common:      db.MakeNewCommon(),
		Token:       uuid.New().String(),
		Name:        req.Name,
		Chats_count: 0,
	}
	err = server.dBWrapper.InsertApplication(a).Error

	if err != nil {
		failure(w, r, http.StatusBadRequest, err.Error())
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to addApplication with error %v", err)
		return
	}

	res := &AddApplicationRes{Application: *a}
	success(w, r, res, logger.ESSENTIAL)
}

func (server *Server) addChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		failure(w, r, http.StatusBadRequest, "invalid request")
		return
	}

	vars := mux.Vars(r)
	appToken, ok := vars["appToken"]
	if !ok {
		failure(w, r, http.StatusBadRequest, "application token is missing")
		return
	}

	logger.LogInfo(logger.SERVER, logger.NON_ESSENTIAL, "The http received message to addChat is %+v", nil)

	//confirm the chat belongs to an app
	a := &db.Application{}
	err := server.dBWrapper.GetApplicationByToken(a, appToken).Error

	if err != nil {
		failure(w, r, http.StatusBadRequest, "application token incorrect")
		return
	}

	c := &db.Chat{
		Common:            db.MakeNewCommon(),
		Application_token: appToken,
		Number:            req.Name,
		Messages_count:    0,
	}

	server.Mq.Publish(q.ENTITIES_QUEUE)

	//success(w, r, res, logger.ESSENTIAL)
}
