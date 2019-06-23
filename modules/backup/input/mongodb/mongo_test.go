package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongoBackup(t *testing.T) {
	d := &BackupInputMongodb{}
	c := GetConfig().(*MongoConfig)
	DSN = "mongodb://127.0.0.1:27017"
	require.NoError(t, InitModule(c))
	names, err := getDatabases()
	assert.NoError(t, err)
	assert.Equal(t, []string{"admin", "local", "test"}, names)
	collections, err := getCollections("test")
	assert.NoError(t, err)
	assert.Equal(t, []string{"link", "test"}, collections)

}
