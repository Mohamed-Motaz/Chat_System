package Server

import (
	db "Server/Database"
	logger "Server/Logger"
	q "Server/MessageQueue"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (server *Server) addChat(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appToken, ok := vars["appToken"]
	if !ok {
		failure(w, r, http.StatusBadRequest, "application token is missing")
		return
	}

	logger.LogInfo(logger.SERVER, logger.NON_ESSENTIAL, "The http received message to addChat is %+v", nil)

	//confirm the chat belongs to an app
	//todo use cache instead
	a := &db.Application{}
	err := server.dBWrapper.GetApplicationByToken(a, appToken).Error

	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Database is down")
		return
	}

	if a.Id == 0 {
		failure(w, r, http.StatusBadRequest, "application token incorrect")
		return
	}

	//get the chat's number from the cache
	//assume its currently highly fault tolerant
	number, err := server.cache.Incr(server.cache.MakeChatCacheKey(appToken))
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Caching layer is down")
		return
	}

	c := &db.Chat{
		Common:            db.MakeNewCommon(),
		Application_token: appToken,
		Number:            int32(number),
		Messages_count:    0,
	}

	toPublish := new(bytes.Buffer)
	err = json.NewEncoder(toPublish).Encode(c)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}
	err = server.Mq.Publish(q.ENTITIES_QUEUE, toPublish.Bytes())

	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		server.cache.Decr(server.cache.MakeChatCacheKey(appToken)) //if this returns an error, then should most probably log it
		return
	}

	res := &AddChatRes{
		Number: c.Number,
	}

	success(w, r, res, logger.ESSENTIAL)
}
