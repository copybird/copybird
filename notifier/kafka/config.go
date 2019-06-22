package kafka

type Config struct {
	MaxRetry   int
	BrokerList []string
	Topic      string
	Message    string
}
