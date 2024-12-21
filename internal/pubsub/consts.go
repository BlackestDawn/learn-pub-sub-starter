package pubsub

type SimpleQueueType int

type AckType int

const (
	QueueTypeDurable SimpleQueueType = iota
	QueueTypeTransient
)

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

const AmqpServer = "amqp://guest:guest@localhost:5672/"
