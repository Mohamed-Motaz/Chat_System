package Server

import (
	db "Server/Database"
	logger "Server/Logger"
	q "Server/MessageQueue"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
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

	//check if the token is indeed in cache
	//if not, call the db once and confirm its presence
	cacheKey := server.cache.MakeChatCacheKey(appToken)
	_, err := server.cache.Get(cacheKey)

	if err != nil {
		if err != redis.Nil {
			failure(w, r, http.StatusInternalServerError, "Caching layer is down")
			return
		}

		//so the key isn't present in cache. Now I need to call the db to confirm its presence
		err := server.confirmAppTokenInDb(w, r, appToken)
		if err != nil {
			return
		}
	}

	//now we can confirm that the token is indeed in the cache. This means that we won't really need to
	//call the db in the future, rather just the cache

	//get the chat's number from the cache
	//assume its currently highly fault tolerant
	number, err := server.cache.Incr(cacheKey)
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

func (server *Server) addMessage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appToken, ok := vars["appToken"]
	if !ok {
		failure(w, r, http.StatusBadRequest, "application token is missing")
		return
	}

	logger.LogInfo(logger.SERVER, logger.NON_ESSENTIAL, "The http received message to addChat is %+v", nil)

	//check if the token is indeed in cache
	//if not, call the db once and confirm its presence
	cacheKey := server.cache.MakeChatCacheKey(appToken)
	_, err := server.cache.Get(cacheKey)

	if err != nil {
		if err != redis.Nil {
			failure(w, r, http.StatusInternalServerError, "Caching layer is down")
			return
		}

		//so the key isn't present in cache. Now I need to call the db to confirm its presence
		err := server.confirmAppTokenInDb(w, r, appToken)
		if err != nil {
			return
		}
	}

	//now we can confirm that the token is indeed in the cache. This means that we won't really need to
	//call the db in the future, rather just the cache

	//get the chat's number from the cache
	//assume its currently highly fault tolerant
	number, err := server.cache.Incr(cacheKey)
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

func (server *Server) confirmAppTokenInDb(w http.ResponseWriter, r *http.Request, appToken string) error {
	logger.LogInfo(logger.DATABASE, logger.NON_ESSENTIAL, "About to confirm this app token %+v", appToken)
	a := &db.Application{}
	err := server.dBWrapper.GetApplicationByToken(a, appToken).Error

	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Database is down")
		return fmt.Errorf("Database is down")
	}

	if a.Id == 0 {
		failure(w, r, http.StatusBadRequest, "application token incorrect")
		return fmt.Errorf("application token incorrect")
	}
	return nil
}
