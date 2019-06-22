package rabbitmq

// Config for RabbitMQ
type Config struct {
	RabbitMQURL   string
	QueueConfig   QueueConfig
	PublishConfig PublishConfig
	MsgConfig     MsgConfig
}

// QueueConfig config parameters for QueueDeclare
type QueueConfig struct {
	Name       string // queue name
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

// PublishConfig config for publish
type PublishConfig struct {
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
}

// MsgConfig config for publish message
type MsgConfig struct {
	ContentType string
	Body        []byte
}
