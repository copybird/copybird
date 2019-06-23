package etcd

import (
	"io"

	"github.com/copybird/copybird/core"
)

// Module Constants
const (
	GroupName  = "backup"
	TypeName   = "input"
	ModuleName = "etcd"
)

type (
	// BackupInputEtcd is struct storing inner properties for mysql backups
	BackupInputEtcd struct {
		core.Module
		reader io.Reader
		writer io.Writer
		config *Config
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
	return nil
}

// Run dumps database
func (b *BackupInputEtcd) Run() error {
	return nil
}

// Close closes ...
func (b *BackupInputEtcd) Close() error {
	return nil
}
