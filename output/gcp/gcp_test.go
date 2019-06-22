package gcp

import (
	"testing"
)

//InitOutput initializes S3 with session
func TestInitOutput(t *testing.T) {

	var gcp GCP

	config := map[string]string{
		"AWS_REGION": "us-east-1",
	}

	err := gcp.InitOutput(config)
	if err == nil {
		t.Error("Should fail to find credentials")
	}
}
