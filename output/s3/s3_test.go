package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//InitOutput initializes S3 with session
func TestInitOutput(t *testing.T) {
	var s S3
	err := s.InitOutput(map[string]string{"AWS_REGION": "us-east-1"})
	assert.NoError(t, err, "should not be any error here")
}
