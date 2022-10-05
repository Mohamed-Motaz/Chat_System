package Server

import (
	db "Server/Database"
	logger "Server/Logger"
	q "Server/MessageQueue"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

	toPublish, err := server.makeChatMqEntity(w, r, q.INSERT_ACTION, c)
	if err != nil {
		return
	}
	err = server.Mq.Publish(q.ENTITIES_QUEUE, toPublish)

	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		server.cache.Decr(server.cache.MakeChatCacheKey(appToken)) //if this returns an error, then should most probably log it
		return
	}

	res := &AddEntityRes{
		Number: c.Number,
	}

	success(w, r, res, logger.ESSENTIAL)
}

func (server *Server) addMessage(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "unable to read request")
		return
	}

	req := &AddMessageReq{}
	err = json.Unmarshal(b, req)
	if err != nil {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to parse request to addMessage %v\nwith error %v", string(b), err)
		failure(w, r, http.StatusBadRequest, "unable to martial request")
		return
	}
	req.Body = strings.TrimSpace(req.Body)
	if len(req.Body) == 0 {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "No body in request to addMessage %v", err)
		failure(w, r, http.StatusBadRequest, "no valid body")
		return
	}

	vars := mux.Vars(r)
	appToken, ok := vars["appToken"]
	if !ok || len(strings.TrimSpace(appToken)) == 0 {
		failure(w, r, http.StatusBadRequest, "application token is missing")
		return
	}

	num, ok := vars["chatNum"]
	if !ok {
		failure(w, r, http.StatusBadRequest, "chat number is missing")
		return
	}
	chatNum, err := strconv.Atoi(num)
	if err != nil {
		failure(w, r, http.StatusBadRequest, "invalid chat number")
		return
	}

	logger.LogInfo(logger.SERVER, logger.NON_ESSENTIAL, "The http received message to addMessage is %+v", nil)

	//check if the token and chat number is indeed in cache
	//if not, call the db once and confirm its presence
	cacheKey := server.cache.MakeMessageCacheKey(appToken, chatNum)
	chatIdCacheKey := server.cache.MakeMessageChatIdCacheKey(appToken, chatNum)

	_, err = server.cache.Get(cacheKey)

	if err != nil {
		if err != redis.Nil {
			failure(w, r, http.StatusInternalServerError, "Caching layer is down")
			return
		}

		//so the key isn't present in cache. Now I need to call the db to confirm its presence
		id, err := server.confirmChatNumberInDb(w, r, appToken, chatNum)
		if err != nil {
			return
		}
		//now i am sure this chat is indeed in the db, so I need to set this chat id in redis for faster access
		err = server.cache.Set(chatIdCacheKey, id, 0)
		if err != nil {
			failure(w, r, http.StatusInternalServerError, "Caching layer is down")
			return
		}
	}

	//now we can confirm that the token and chat number is indeed in the cache. This means that we won't really need to
	//call the db in the future, rather just the cache

	number, err := server.cache.Incr(cacheKey)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Caching layer is down")
		return
	}

	chatId, err := server.cache.Get(chatIdCacheKey)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Caching layer is down")
		return
	}

	cId, err := strconv.Atoi(chatId)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Issue while retrieving data from cache")
		return
	}

	m := &db.Message{
		Common:  db.MakeNewCommon(),
		Chat_id: int32(cId),
		Number:  int32(number),
		Body:    req.Body,
	}

	toPublish, err := server.makeMessageMqEntity(w, r, q.INSERT_ACTION, m)
	if err != nil {
		return
	}
	err = server.Mq.Publish(q.ENTITIES_QUEUE, toPublish)

	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		server.cache.Decr(server.cache.MakeChatCacheKey(appToken)) //if this returns an error, then should most probably log it
		return
	}

	res := &AddEntityRes{
		Number: m.Number,
	}

	success(w, r, res, logger.ESSENTIAL)
}

//confirm the app token is indeed in the db
func (server *Server) confirmAppTokenInDb(w http.ResponseWriter, r *http.Request, appToken string) error {
	logger.LogInfo(logger.DATABASE, logger.NON_ESSENTIAL, "About to confirm this app token %+v", appToken)
	a := &db.Application{}
	err := server.dBWrapper.GetApplicationByToken(a, appToken).Error

	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Database is down")
		return fmt.Errorf("database is down")
	}

	if a.Id == 0 {
		failure(w, r, http.StatusBadRequest, "application token incorrect")
		return fmt.Errorf("application token incorrect")
	}
	return nil
}

//confirm the chat number is indeed in the db and return its id
func (server *Server) confirmChatNumberInDb(w http.ResponseWriter, r *http.Request, appToken string, chatNum int) (int, error) {
	logger.LogInfo(logger.DATABASE, logger.NON_ESSENTIAL, "About to confirm this app token %+v", appToken)
	id := 0
	err := server.dBWrapper.GetChatIdByAppTokenAndNumber(&id, appToken, chatNum).Error

	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Database is down")
		return 0, fmt.Errorf("database is down")
	}

	if id == 0 {
		failure(w, r, http.StatusBadRequest, "application token or chat number incorrect")
		return 0, fmt.Errorf("application token or chat incorrect")
	}
	return id, nil
}

func (server *Server) makeChatMqEntity(w http.ResponseWriter, r *http.Request, action q.DB_ACTION, chat *db.Chat) ([]byte, error) {
	c := q.Chat{
		Id:                chat.Id,
		Application_token: chat.Application_token,
		Number:            chat.Number,
		Messages_count:    chat.Messages_count,
		Created_at:        chat.Created_at,
		Updated_at:        chat.Updated_at,
	}

	cBytes := new(bytes.Buffer)
	err := json.NewEncoder(cBytes).Encode(c)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		return nil, err
	}

	obj := q.TransferObj{
		Action:  action,
		ObjType: q.CHAT,
		Bytes:   cBytes.Bytes(),
	}

	resBytes := new(bytes.Buffer)
	err = json.NewEncoder(resBytes).Encode(obj)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		return nil, err
	}

	return resBytes.Bytes(), nil
}

func (server *Server) makeMessageMqEntity(w http.ResponseWriter, r *http.Request, action q.DB_ACTION, message *db.Message) ([]byte, error) {
	m := q.Message{
		Id:         message.Id,
		Chat_id:    message.Chat_id,
		Number:     message.Number,
		Body:       message.Body,
		Created_at: message.Created_at,
		Updated_at: message.Updated_at,
	}

	mBytes := new(bytes.Buffer)
	err := json.NewEncoder(mBytes).Encode(m)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		return nil, err
	}

	obj := q.TransferObj{
		Action:  action,
		ObjType: q.MESSAGE,
		Bytes:   mBytes.Bytes(),
	}

	resBytes := new(bytes.Buffer)
	err = json.NewEncoder(resBytes).Encode(obj)
	if err != nil {
		failure(w, r, http.StatusInternalServerError, "Internal server error")
		return nil, err
	}

	return resBytes.Bytes(), nil
}
