package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKafka(t *testing.T) {
	k := &GlobalNotifieKafka{}
	c := GetConfig().(*Config)
	assert.NotNil(t, c)
	BrokerList = []string{"localhost:9092", "localhost:9092"}
	Topic = "hello"
	Message = "world"
	assert.NoError(t, InitModule(c))
	assert.NoError(t, Run())
}
