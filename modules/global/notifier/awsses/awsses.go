package awsses

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/copybird/copybird/core"
)

const (
	GROUP_NAME  = "global"
	TYPE_NAME   = "notifier"
	MODULE_NAME = "awsses"
)

type GlobalNotifierAwsses struct {
	core.Module
	config  *Config
	simem   *ses.SES // simple email service
	seInput *ses.SendEmailInput
	reader  io.Reader
	writer  io.Writer
}

func (m *GlobalNotifierAwsses) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierAwsses) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *GlobalNotifierAwsses) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierAwsses) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierAwsses) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierAwsses) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)

	// Create a new session in the config region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(m.config.Region)},
	)
	if err != nil {
		return err
	}

	// Create an SES session.
	svc := ses.New(sess)
	m.simem = svc

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(m.config.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(m.config.Charset),
					Data:    aws.String(m.config.HTMLbody),
				},
				Text: &ses.Content{
					Charset: aws.String(m.config.Charset),
					Data:    aws.String(m.config.Textbody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(m.config.Charset),
				Data:    aws.String(m.config.Subject),
			},
		},
		Source: aws.String(m.config.Sender),
		// TODO: Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}
	m.seInput = input

	return nil
}

func (m *GlobalNotifierAwsses) Run(ctx context.Context) error {

	// Attempt to send the email.
	_, err := m.simem.SendEmail(m.seInput)

	// Display error messages if they occur.
	if err != nil {
		return err
	}

	return nil
}

// Close closes compressor
func (m *GlobalNotifierAwsses) Close() error {
	return nil
}
