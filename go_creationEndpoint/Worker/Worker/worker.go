package Worker

import (
	db "Worker/Database"
	logger "Worker/Logger"
	q "Worker/MessageQueue"
	"encoding/json"
	"fmt"
	"time"
)

//create a new Worker
func NewWorker() (*Worker, error) {

	worker := &Worker{
		dBWrapper: db.New(),
		Mq:        q.New("amqp://" + q.MqUsername + ":" + q.MqPassword + "@" + q.MqHost + ":" + q.MqPort + "/"),
	}

	go worker.qConsumer()

	return worker, nil
}

//
// start a thread that waits on the entities from the message queue
//
func (worker *Worker) qConsumer() {
	ch, err := worker.Mq.Consume(q.ENTITIES_QUEUE)

	if err != nil {
		logger.FailOnError(logger.WORKER, logger.ESSENTIAL, "Worker can't consume entities because with this error %v", err)
	}

	for {
		select {
		case consumed := <-ch: //entity has been pushed to the queue
			body := consumed.Body

			data := &q.TransferObj{}

			err := json.Unmarshal(body, data)
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to unmarshal entity %v with error %v\nWill discard it", string(body), err)
				consumed.Ack(false)
				continue
			}

			logger.LogInfo(logger.WORKER, logger.ESSENTIAL, "Received this entity from %+v from the queue", data)

			worker.doWork(data)

			consumed.Ack(false) //ack after everything is done. This should be blocking other workers which is definetly a bottleneck
		default:
			time.Sleep(time.Second * 2)
		}
	}
}

func (worker *Worker) doWork(obj *q.TransferObj) {
	logger.LogInfo(logger.WORKER, logger.ESSENTIAL, "The transfer obj received %+v", obj)
	if obj.Action == q.INSERT_ACTION {
		switch obj.ObjType {
		case q.CHAT:
			c := &q.Chat{}
			err := json.Unmarshal(obj.Bytes, c)
			fmt.Printf("This is the data received %+v\n", c)
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to parse request to insert chart %v\nwith error %v", string(obj.Bytes), err)
				return
			}
			chat := makeDbChatFromQChat(c)
			err = worker.dBWrapper.InsertChat(chat).Error
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to insert chat with error %v", err)
				return
			}
		case q.MESSAGE:
			m := &q.Message{}
			err := json.Unmarshal(obj.Bytes, m)
			fmt.Printf("This is the data received %+v\n", m)
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to parse request to insert message %v\nwith error %v", string(obj.Bytes), err)
				return
			}
			message := makeDbMessageFromQMessage(m)
			err = worker.dBWrapper.InsertMessage(message).Error
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to insert message with error %v", err)
				return
			}
		default:
			logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to decode obj from the queue %+v", obj)
		}

	} else if obj.Action == q.UPDATE_ACTION {
		switch obj.ObjType {
		case q.MESSAGE:
			m := &q.Message{}
			err := json.Unmarshal(obj.Bytes, m)
			fmt.Printf("This is the data received %+v\n", m)
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to parse request to update message %v\nwith error %v", string(obj.Bytes), err)
				return
			}
			message := makeDbMessageFromQMessage(m)
			err = worker.dBWrapper.UpdateMessage(message).Error
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to update message with error %v", err)
				return
			}
		case q.CHATS_CTR:
			chatsCtr := &q.ChatsCtr{}
			err := json.Unmarshal(obj.Bytes, chatsCtr)
			fmt.Printf("This is the data received %+v\n", chatsCtr)
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to parse request to update the chats ctr %v\nwith error %v", string(obj.Bytes), err)
				return
			}
			err = worker.dBWrapper.UpdateApplicationCtr(chatsCtr.Token, int(chatsCtr.Chats_count), time.Now()).Error
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to update the chats ctr with error %v", err)
				return
			}
		case q.MESSAGES_CTR:
			messageCtr := &q.MessagesCtr{}
			err := json.Unmarshal(obj.Bytes, messageCtr)
			fmt.Printf("This is the data received %+v\n", messageCtr)
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to parse request to update the messages ctr %v\nwith error %v", string(obj.Bytes), err)
				return
			}
			err = worker.dBWrapper.UpdateChatsCtr(messageCtr.Application_token, int(messageCtr.Number), int(messageCtr.Messages_count), time.Now()).Error
			if err != nil {
				logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to update the messages ctr with error %v", err)
				return
			}
		default:
			logger.LogError(logger.WORKER, logger.ESSENTIAL, "Unable to decode obj from the queue %+v", obj)
		}
	}
}

//converts from q.Chat to db.Chat
func makeDbChatFromQChat(chat *q.Chat) *db.Chat {
	return &db.Chat{
		Common: db.Common{
			Id:         chat.Id,
			Created_at: chat.Created_at,
			Updated_at: chat.Updated_at,
		},
		Application_token: chat.Application_token,
		Number:            chat.Number,
		Messages_count:    chat.Messages_count,
	}
}

//converts from q.Message to db.Message
func makeDbMessageFromQMessage(message *q.Message) *db.Message {
	return &db.Message{
		Common: db.Common{
			Id:         message.Id,
			Created_at: message.Created_at,
			Updated_at: message.Updated_at,
		},
		Chat_id: message.Chat_id,
		Number:  message.Number,
		Body:    message.Body,
	}
}
