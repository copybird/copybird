package awssqs

type Config struct {
	Region                 string
	AccountAccessKeyID     string
	AccountSecretAccessKey string
	Queues                 string
	MessageBody            string
}
