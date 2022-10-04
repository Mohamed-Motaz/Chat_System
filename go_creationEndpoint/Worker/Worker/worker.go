package Worker

import (
	db "Worker/Database"
	q "Worker/MessageQueue"
)

//create a new Worker
func NewWorker() (*Worker, error) {

	worker := &Worker{
		dBWrapper: db.New(),
		Mq:        q.New("amqp://" + q.MqUsername + ":" + q.MqPassword + "@" + q.MqHost + ":" + q.MqPort + "/"),
	}

	return worker, nil
}
