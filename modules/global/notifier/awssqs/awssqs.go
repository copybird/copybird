package awssqs

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/copybird/copybird/core"
	"github.com/knative/pkg/cloudevents"
)

const (
	MODULE_NAME = "awssqs"
	GROUP_NAME  = "global"
	TYPE_NAME   = "notifier"
)

func (m *GlobalNotifierAWSSQS) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierAWSSQS) GetType() core.ModuleType {
	return TYPE_NAME
}

type GlobalNotifierAWSSQS struct {
	core.Module
	config *Config
	reader io.Reader
	writer io.Writer
}

func (m *GlobalNotifierAWSSQS) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierAWSSQS) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierAWSSQS) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	return nil
}

func (m *GlobalNotifierAWSSQS) Run() error {
	if err := m.config.NotifyAWSSQS(); err != nil {
		return err
	}

	return nil
}

func (m *GlobalNotifierAWSSQS) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierAWSSQS) Close() error {
	return nil
}

type Clients struct {
	SQS         sqsiface.SQSAPI
	CloudEvents *cloudevents.Client
}

func (c *Config) NotifyAWSSQS() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccountAccessKeyID, c.AccountSecretAccessKey, ""),
		MaxRetries:  aws.Int(5),
	})

	if err != nil {
		return err
	}

	sqsClient := sqs.New(sess)

	queueUrls, err := sqsClient.ListQueues(&sqs.ListQueuesInput{QueueNamePrefix: aws.String(c.Queues)})

	if err != nil {
		return err
	}

	sendMessage := &sqs.SendMessageInput{
		MessageBody:  aws.String(c.MessageBody),
		QueueUrl:     queueUrls.QueueUrls[0],
		DelaySeconds: aws.Int64(3),
	}

	_, err = sqsClient.SendMessage(sendMessage)
	if err != nil {
		return err
	}

	return nil
}
