package mongodb

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongoBackup2(t *testing.T) {
	d := &BackupInputMongodb{}
	c := d.GetConfig().(*MongoConfig)
	c.DSN = "mongodb://127.0.0.1:27017"
	require.NoError(t, d.InitModule(c))
	names, err := d.getDatabases(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, []string{"admin", "local", "test"}, names)
	collections, err := d.getCollections(context.TODO(), "test")
	assert.NoError(t, err)
	assert.Equal(t, []string{"link", "test"}, collections)

}

func TestExportCollection(t *testing.T) {
	d := &BackupInputMongodb{}
	c := d.GetConfig().(*MongoConfig)
	c.DSN = "mongodb://127.0.0.1:27017"
	require.NoError(t, d.InitModule(c))
	tmpFile, err := os.Create("./export.txt")
	if !assert.NoError(t, err, "tmp file") {
		t.Fail()
	}

	d.writer = tmpFile

	assert.NoError(t, d.exportCollection("admin", "system.version"), "export collection")
	tmpFile.Close()
	// os.Remove(tmpFile.Name())
}
