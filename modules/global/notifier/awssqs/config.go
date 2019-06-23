package awssqs

type Config struct {
	Region                 string
	AccountAccessKeyID     string
	AccountSecretAccessKey string
	QueueUrl               string
	MessageBody            string
}
