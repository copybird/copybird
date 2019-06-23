package tar

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/copybird/copybird/core"
	"github.com/davecgh/go-spew/spew"
)

// Module Constants
const (
	GroupName  = "backup"
	TypeName   = "input"
	ModuleName = "tar"
)

type (
	// BackupInputTar is struct storing inner properties for mysql backups
	BackupInputTar struct {
		core.Module
		reader io.Reader
		writer io.Writer
		config *Config
	}
)

// GetGroup returns group
func (b *BackupInputTar) GetGroup() core.ModuleGroup {
	return GroupName
}

// GetType returns type
func (b *BackupInputTar) GetType() core.ModuleType {
	return TypeName
}

// GetName returns name of module
func (b *BackupInputTar) GetName() string {
	return ModuleName
}

// GetConfig returns config of module
func (b *BackupInputTar) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (b *BackupInputTar) InitPipe(w io.Writer, r io.Reader) error {
	b.reader = r
	b.writer = w
	return nil
}

// InitModule initializes module
func (b *BackupInputTar) InitModule(cfg interface{}) error {
	b.config = cfg.(*Config)
	return nil
}

// Run dumps database
func (b *BackupInputTar) Run() error {
	if _, err := os.Stat(b.config.DirectoryPath); err != nil {
		return err
	}
	tw := tar.NewWriter(b.writer)
	defer tw.Close()
	return filepath.Walk(b.config.DirectoryPath, func(file string, fi os.FileInfo, err error) error {
		spew.Dump(file)

		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, b.config.DirectoryPath, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})
	return nil
}

// Close closes ...
func (b *BackupInputTar) Close() error {
	return nil
}
