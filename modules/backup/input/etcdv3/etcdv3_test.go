package etcdv3

import (
	"bytes"
	"context"
	"testing"

	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/assert"
)

func TestEtcdBackup(t *testing.T) {
	b := &BackupInputEtcd{}
	assert.Equal(t, &Config{}, b.GetConfig())
	conf := &Config{Endpoints: []string{"http://127.0.0.1:23791"}, Key: "/test"}
	assert.NoError(t, b.InitModule(conf))
	wr := &bytes.Buffer{}
	assert.NoError(t, b.InitPipe(wr, nil))
	c, err := clientv3.New(clientv3.Config{Endpoints: conf.Endpoints})
	assert.NoError(t, err)
	a := clientv3.NewKV(c)
	a.Put(context.TODO(), conf.Key, "")
	a.Put(context.TODO(), conf.Key+"/test", "World")
	a.Put(context.TODO(), conf.Key+"/dir", "Hey")
	a.Put(context.TODO(), conf.Key+"/dir/name", "Value")
	assert.NoError(t, b.Run(context.TODO()))
	a.Delete(context.TODO(), conf.Key, clientv3.WithPrefix())

}
