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
		{"c0pybird0@gmail.com", "c0pybird-admin", "omuraliev.baurzhan@gmail.com"},
	}

	for _, tc := range testCase {
		conf := Email{Config: &Config{MailerUser: tc.MailerUser, MailerPassword: tc.MailerPassword, MailTo: tc.MailTo}}
		err := SendEmail()
		assert.Equal(t, nil, err)
	}
}
