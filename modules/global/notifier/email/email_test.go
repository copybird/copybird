package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {

	testCase := []struct {
		MailerUser     string
		MailerPassword string
		MailTo         string
	}{
		{"c0pybird0@gmail.com", "pas$$w0rd", "example.com"},
	}

	for _, tc := range testCase {
		g := &GlobalNotifierEmail{}
		assert.NotNil(t, g.GetConfig())
		assert.NoError(t, g.InitModule(&Config{MailerUser: tc.MailerUser, MailerPassword: tc.MailerPassword, MailTo: tc.MailTo}))
		assert.Error(t, g.SendEmail())
	}
}
