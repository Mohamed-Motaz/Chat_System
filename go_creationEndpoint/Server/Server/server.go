package Server

import (
	db "Server/Database"
	logger "Server/Logger"
	q "Server/MessageQueue"
	c "Server/RedisCache"
	utils "Server/Utils"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

//create a new server
func NewServer() (*Server, error) {

	server := &Server{
		dBWrapper: db.New(),
		Mq:        q.New("amqp://" + q.MqUsername + ":" + q.MqPassword + "@" + q.MqHost + ":" + q.MqPort + "/"),
		cache:     c.New(c.CacheHost + ":" + c.CachePort),
	}
	r := mux.NewRouter()
	registerRoutes(r, server)

	go server.serve(r)
	go server.sendNumberUpdatesToWorker()

	return server, nil
}

func registerRoutes(r *mux.Router, s *Server) {
	r.HandleFunc("/api/v1/applications/{appToken}/chats", s.addChat).Methods("POST")
	r.HandleFunc("/api/v1/applications/{appToken}/chats/{chatNum}/messages", s.addMessage).Methods("POST")
	r.HandleFunc("/api/v1/applications/{appToken}/chats/{chatNum}/messages/{messageNum}", s.updateMessage).Methods("PUT")

}

func (server *Server) serve(r *mux.Router) {
	logger.LogInfo(logger.SERVER, logger.ESSENTIAL, "Listening on %v:%v", MyHost, MyPort)
	err := http.ListenAndServe(MyHost+":"+MyPort, r)

	if err != nil {
		logger.FailOnError(logger.SERVER, logger.ESSENTIAL, "Failed in listening on port with error %v", err)
	}
}

func (server *Server) sendNumberUpdatesToWorker() {

	//map for keeping track of the old values, so as not to keep sending repeated counts
	//to the worker

	mp := make(map[string]int) //this map could be a potential leak, since it only grows
	//a potential solution would be to remove it and create a new one periodically

	for {
		logger.LogInfo(logger.SERVER, logger.ESSENTIAL, "About to send number updates to the server")

		//get all chat ctrs for the different app tokens
		iter := server.cache.GetKeysWithPatternIterator("*.CHAT_CTR")
		for iter.Next(server.cache.GetCacheCtx()) {
			key := iter.Val()
			val, err := server.cache.Get(key)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Caching service error %v", err)
				continue
			}

			count, err := strconv.Atoi(val)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to convert count %v to int with err %v", count, err)
				continue
			}

			token := utils.GetStringInBetween(key, "APP_TOKEN:", ".CHAT_CTR")
			// fmt.Printf("key %+v -- token %+v\n", key, token)

			chatsCtr := q.ChatsCtr{
				Token:       token,
				Chats_count: int32(count),
			}
			prevCount, ok := mp[token]
			if ok && prevCount == int(chatsCtr.Chats_count) {
				//no need to resend the same old counter
				continue
			}
			mp[token] = int(chatsCtr.Chats_count)

			toPublish, err := server.makeMessageMqEntityNoHandler(q.UPDATE_ACTION, q.CHATS_CTR, chatsCtr)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to makeMessageMqEntityNoHandler %+v to bytes with err %v", chatsCtr, err)
				continue
			}

			err = server.Mq.Publish(q.ENTITIES_QUEUE, toPublish)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to publish to mq %+v to bytes with err %v", string(toPublish), err)
				continue
			}
		}
		if err := iter.Err(); err != nil {
			panic(err)
		}

		//get all message ctrs for the different chats
		iter = server.cache.GetKeysWithPatternIterator("*.MESSAGE_CTR")
		for iter.Next(server.cache.GetCacheCtx()) {
			key := iter.Val()
			val, err := server.cache.Get(key)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Caching service error %v", err)
				continue
			}
			count, err := strconv.Atoi(val)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Caching service error %v", err)
				continue
			}

			token := utils.GetStringInBetween(key, "APP_TOKEN:", "--CHAT_NUM")
			num := utils.GetStringInBetween(key, "CHAT_NUM:", ".MESSAGE_CTR")
			number, err := strconv.Atoi(num)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Caching service error %v", err)
				continue
			}

			fmt.Printf("key %+v -- token %+v -- num %v\n", key, token, num)

			msgCtr := q.MessagesCtr{
				Application_token: token,
				Number:            int32(number),
				Messages_count:    int32(count),
			}
			mpKey := fmt.Sprintf("%s-%d", msgCtr.Application_token, msgCtr.Number)
			prevCount, ok := mp[mpKey]
			if ok && prevCount == int(msgCtr.Messages_count) {
				//no need to resend the same old counter
				continue
			}
			mp[mpKey] = int(msgCtr.Messages_count)

			toPublish, err := server.makeMessageMqEntityNoHandler(q.UPDATE_ACTION, q.MESSAGES_CTR, msgCtr)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to makeMessageMqEntityNoHandler %+v to bytes with err %v", msgCtr, err)
				continue
			}

			err = server.Mq.Publish(q.ENTITIES_QUEUE, toPublish)
			if err != nil {
				logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to publish to mq %+v to bytes with err %v", string(toPublish), err)
				continue
			}

		}
		if err := iter.Err(); err != nil {
			panic(err)
		}

		time.Sleep(10 * time.Second)
	}

}

func (server *Server) makeMessageMqEntityNoHandler(action q.DB_ACTION, objType q.OBJ_TYPE, entity interface{}) ([]byte, error) {

	mBytes := new(bytes.Buffer)
	err := json.NewEncoder(mBytes).Encode(entity)
	if err != nil {
		return nil, err
	}

	obj := q.TransferObj{
		Action:  action,
		ObjType: objType,
		Bytes:   mBytes.Bytes(),
	}

	resBytes := new(bytes.Buffer)
	err = json.NewEncoder(resBytes).Encode(obj)
	if err != nil {
		return nil, err
	}

	return resBytes.Bytes(), nil
}
