package s3

import (
	"testing"
)

//InitOutput initializes S3 with session
func TestInitOutput(t *testing.T) {

	var s S3

	config := map[string]string{
		"AWS_REGION": "us-east-1",
	}

	err := s.InitOutput(config)
	if err != nil {
		t.Error(err)
	}
}
