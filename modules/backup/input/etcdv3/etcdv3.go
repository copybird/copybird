package etcdv3

import (
	"context"
	"encoding/json"
	"io"

	"github.com/copybird/copybird/core"
	"github.com/coreos/etcd/clientv3"
	"github.com/davecgh/go-spew/spew"
)

// Module Constants
const (
	GroupName  = "backup"
	TypeName   = "input"
	ModuleName = "etcdv3"
)

type (
	// BackupInputEtcd is struct storing inner properties for mysql backups
	BackupInputEtcd struct {
		core.Module
		reader io.Reader
		writer io.Writer
		config *Config
		api    clientv3.KV
	}
)

// GetGroup returns group
func (b *BackupInputEtcd) GetGroup() core.ModuleGroup {
	return GroupName
}

// GetType returns type
func (b *BackupInputEtcd) GetType() core.ModuleType {
	return TypeName
}

// GetName returns name of module
func (b *BackupInputEtcd) GetName() string {
	return ModuleName
}

// GetConfig returns config of module
func (b *BackupInputEtcd) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (b *BackupInputEtcd) InitPipe(w io.Writer, r io.Reader) error {
	b.reader = r
	b.writer = w
	return nil
}

// InitModule initializes module
func (b *BackupInputEtcd) InitModule(cfg interface{}) error {
	b.config = cfg.(*Config)
	c, err := clientv3.New(clientv3.Config{
		Endpoints: b.config.Endpoints,
	})
	if err != nil {
		return err
	}
	b.api = clientv3.NewKV(c)
	return nil
}

// Run dumps database
func (b *BackupInputEtcd) Run() error {
	resp, err := b.api.Get(context.TODO(), b.config.Key)
	if err != nil {
		return err
	}
	spew.Dump(resp.Kvs)
	j := json.NewEncoder(b.writer)
	return j.Encode(resp.Kvs)
}

// Close closes ...
func (b *BackupInputEtcd) Close() error {
	return nil
}
