package awsses

import (
	"testing"

	"gotest.tools/assert"
)

var (
	noCredErr = "NoCredentialProviders: no valid providers in chain. Deprecated.\n\tFor verbose messaging see aws.Config.CredentialsChainVerboseErrors"
)

func TestAwsSes_NoCredentialProvErrs(t *testing.T) {
	conf := &Config{
		Region: "us-west-2",
	}

	as := &GlobalNotifierAwsses{}
	assert.Assert(t, as.GetConfig() != nil)
	assert.NilError(t, as.InitModule(conf))
	err := as.Run()
	assert.Error(t, err, noCredErr)
}

func TestAwsSes_WithCredential(t *testing.T) {
	conf := &Config{
		Region:    "us-west-2",
		Sender:    "sender@example.com",
		Recipient: "recipient@example.com",
		Subject:   "Amazon SES Test (AWS SDK for Go)",
		HTMLbody:  "Test",
		Textbody:  "This email was sent with Amazon SES using the AWS SDK for Go.",
		Charset:   "UTF-8",
	}

	as := &GlobalNotifierAwsses{}
	assert.Assert(t, as.GetConfig() != nil)
	assert.NilError(t, as.InitModule(conf))
	err := as.Run()
	assert.Error(t, err, noCredErr)
}
