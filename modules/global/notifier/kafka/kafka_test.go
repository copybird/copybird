package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKafka(t *testing.T) {
	k := &GlobalNotifieKafka{}
	c := k.GetConfig().(*Config)
	assert.NotNil(t, c)
	c.BrokerList = []string{"localhost:9092", "localhost:9092"}
	c.Topic = "hello"
	c.Message = "world"
	assert.NoError(t, k.InitModule(c))
	assert.NoError(t, k.Run())
}
