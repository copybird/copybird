package consul

import (
	"io"

	"github.com/copybird/copybird/core"
)

// Module Constants
const (
	GroupName  = "backup"
	TypeName   = "input"
	ModuleName = "consul"
)

type (
	// BackupInputConsul is struct storing inner properties for mysql backups
	BackupInputConsul struct {
		core.Module
		reader io.Reader
		writer io.Writer
		config *Config
	}
)

// GetGroup returns group
func (b *BackupInputConsul) GetGroup() core.ModuleGroup {
	return GroupName
}

// GetType returns type
func (b *BackupInputConsul) GetType() core.ModuleType {
	return TypeName
}

// GetName returns name of module
func (b *BackupInputConsul) GetName() string {
	return ModuleName
}

// GetConfig returns config of module
func (b *BackupInputConsul) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (b *BackupInputConsul) InitPipe(w io.Writer, r io.Reader) error {
	b.reader = r
	b.writer = w
	return nil
}

// InitModule initializes module
func (b *BackupInputConsul) InitModule(cfg interface{}) error {
	b.config = cfg.(*Config)
	return nil
}

// Run dumps database
func (b *BackupInputConsul) Run(ctx context.Context) error {
	return nil
}

// Close closes ...
func (b *BackupInputConsul) Close() error {
	return nil
}
