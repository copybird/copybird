package mongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMongoBackup(t *testing.T) {
	d := &MongoDumper{}
	c := d.GetConfig().(*MongoConfig)
	c.Host = "127.0.0.1"
	require.NoError(t, d.InitModule(c))
}
