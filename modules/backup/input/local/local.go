package local

import (
	"io"
	"os"

	"github.com/copybird/copybird/core"
)

// Module Constants
const (
	GroupName  = "backup"
	TypeName   = "input"
	ModuleName = "local"
)

type (
	// BackupInputLocal is struct storing inner properties for mysql backups
	BackupInputLocal struct {
		core.Module
		reader io.Reader
		writer io.Writer
		config *Config
	}
)

// GetGroup returns group
func (b *BackupInputLocal) GetGroup() core.ModuleGroup {
	return GroupName
}

// GetType returns type
func (b *BackupInputLocal) GetType() core.ModuleType {
	return TypeName
}

// GetName returns name of module
func (b *BackupInputLocal) GetName() string {
	return ModuleName
}

// GetConfig returns config of module
func (b *BackupInputLocal) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (b *BackupInputLocal) InitPipe(w io.Writer, r io.Reader) error {
	b.reader = r
	b.writer = w
	return nil
}

// InitModule initializes module
func (b *BackupInputLocal) InitModule(cfg interface{}) error {
	b.config = cfg.(*Config)
	return nil
}

// Run dumps database
func (b *BackupInputLocal) Run() error {
	f, err := os.Open(b.config.Filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(b.writer, f); err != nil {
		return err
	}
	return nil
}

// Close closes ...
func (b *BackupInputLocal) Close() error {
	return nil
}
