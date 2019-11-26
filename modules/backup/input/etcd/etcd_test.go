// +build external

package etcd

import (
	"bytes"
	"context"
	"testing"

	"github.com/etcd-io/etcd/client"
	"github.com/stretchr/testify/assert"
)

func TestEtcdBackup(t *testing.T) {
	b := &BackupInputEtcd{}
	assert.Equal(t, &Config{}, b.GetConfig())
	conf := &Config{Endpoints: []string{"http://127.0.0.1:23791"}, Key: "/test"}
	assert.NoError(t, b.InitModule(conf))
	wr := &bytes.Buffer{}
	assert.NoError(t, b.InitPipe(wr, nil))
	c, err := client.New(client.Config{Endpoints: conf.Endpoints, Transport: client.DefaultTransport})
	assert.NoError(t, err)
	a := client.NewKeysAPI(c)
	o := client.SetOptions{Dir: true}
	a.Set(context.TODO(), conf.Key, "", &o)
	a.Set(context.TODO(), conf.Key+"/test", "World", nil)
	a.Set(context.TODO(), conf.Key+"/dir", "Hey", &o)
	a.Set(context.TODO(), conf.Key+"/dir/name", "Value", nil)
	assert.NoError(t, b.Run(context.TODO()))
	a.Delete(context.TODO(), conf.Key, &client.DeleteOptions{Recursive: true})

}
