package MessageQueue

import (
	logger "Worker/Logger"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New(amqpAddr string) *MQ {
	mq := &MQ{
		conn: nil,
		ch:   nil,
		qMap: make(map[string]*amqp.Queue),
		mu:   sync.Mutex{},
	}

	//try to reconnect and sleep 2 seconds on failure
	var err error = fmt.Errorf("error")
	ctr := 0

	for ctr < 3 && err != nil {
		err = mq.connect(amqpAddr)
		ctr++
		if err != nil {
			time.Sleep(time.Second * 10)
		}
	}

	if err != nil {
		logger.FailOnError(logger.DATABASE, logger.ESSENTIAL, "Unable to setup the MessageQueue")
	} else {
		logger.LogInfo(logger.DATABASE, logger.ESSENTIAL, "MessageQueue setup complete")
	}
	return mq
}

func (mq *MQ) Close() {
	mq.ch.Close()
	mq.conn.Close()
}

func (mq *MQ) Consume(qName string) (<-chan amqp.Delivery, error) {

	//if q is nill, declare it
	mq.mu.Lock()
	if _, ok := mq.qMap[qName]; !ok {
		mq.mu.Unlock()
		q, err := mq.ch.QueueDeclare(
			qName,
			true,  //durable
			false, //autoDelete
			false, //exclusive -> send errors when another consumer tries to connect to it
			false, //noWait
			nil,
		)
		if err != nil {
			return nil, err
		}
		mq.mu.Lock()
		mq.qMap[qName] = &q
	}
	mq.mu.Unlock()

	err := mq.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}
	msgs, err := mq.ch.Consume(
		qName,
		"",    //consumer --> unique identifier that rabbit can just generate
		false, //auto ack
		false, //exclusive --> errors if other consumers consume
		false, //nolocal, dont receive back if i send on channel
		false, //no wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	logger.LogInfo(logger.MESSAGE_Q, logger.NON_ESSENTIAL, "Successfully subscribed to queue %v", qName)
	return msgs, nil
}

func (mq *MQ) connect(amqpAddr string) error {

	conn, err := amqp.Dial(amqpAddr)
	if err != nil {
		logger.FailOnError(logger.MESSAGE_Q, logger.ESSENTIAL, "Exiting because of err while establishing connection with message queue %v", err)
	}
	mq.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		logger.FailOnError(logger.MESSAGE_Q, logger.ESSENTIAL, "Exiting because of err while opening a channel %v", err)
	}
	mq.ch = ch

	logger.LogInfo(logger.MESSAGE_Q, logger.ESSENTIAL, "Connection to message queue has been successfully established")
	return nil
}
