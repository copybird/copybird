package awsses

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	MODULE_NAME = "awsses"
)

type AwsSes struct {
	config  *Config
	simem   *ses.SES //simple email service
	seInput *ses.SendEmailInput
	reader  io.Reader
	writer  io.Writer
}

func (n *AwsSes) GetName() string {
	return MODULE_NAME
}

func (n *AwsSes) GetConfig() interface{} {
	return &Config{}
}

func (c *AwsSes) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *AwsSes) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)
	c.config = cfg

	// Create a new session in the config region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.config.Region)},
	)
	if err != nil {
		return err
	}

	// Create an SES session.
	svc := ses.New(sess)
	c.simem = svc

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(c.config.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(c.config.CharSet),
					Data:    aws.String(c.config.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(c.config.CharSet),
					Data:    aws.String(c.config.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(c.config.CharSet),
				Data:    aws.String(c.config.Subject),
			},
		},
		Source: aws.String(c.config.Sender),
		// TODO: Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}
	c.seInput = input

	return nil
}

func (c *AwsSes) Run() error {

	// Attempt to send the email.
	_, err := c.simem.SendEmail(c.seInput)

	// Display error messages if they occur.
	if err != nil {
		return err
	}

	return nil
}

// Close closes compressor
func (c *AwsSes) Close() error {
	return nil
}
