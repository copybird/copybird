package rabbitmq

// Config for GlobalNotifierRabbitmq
type Config struct {
	RabbitMQURL      string
	QueueName        string
	QueueDurable     bool
	QueueAutoDelete  bool
	QueueExclusive   bool
	QueueNoWait      bool
	PublishExchange  string
	PublishKey       string
	PublishMandatory bool
	PublishImmediate bool
	MsgContentType   string
	MsgBody          string
}
