package common

import (
	"fmt"
	compress_gzip "github.com/copybird/copybird/compress/gzip"
	"github.com/copybird/copybird/core"
	"github.com/copybird/copybird/encryption/aesgcm"
	"github.com/copybird/copybird/input/mysql"
	"github.com/copybird/copybird/output/local"
	"log"

	// lz4_compress "github.com/copybird/copybird/compress/lz4"
	// lz4_decompress "github.com/copybird/copybird/decompress/lz4"
)

type ModuleType int

const (
	ModuleTypeConnect ModuleType = iota
	ModuleTypeInput
	ModuleTypeCompress
	ModuleTypeEncryption
	ModuleTypeOutput
	ModuleTypeNotify
)

func (m ModuleType) String() string {
	return [...]string{
		"connect",
		"input",
		"compress",
		"encryption",
		"output",
		"notify",
	}[m]
}

func (a *App) Setup() error {
	a.RegisterModules()
	return nil
}

func (a *App) RegisterModules() {
	a.RegisterModule(ModuleTypeInput, &mysql.MySQLDumper{})
	a.RegisterModule(ModuleTypeCompress, &compress_gzip.Compress{})
	a.RegisterModule(ModuleTypeEncryption, &aesgcm.EncryptionAESGCM{})
	a.RegisterModule(ModuleTypeOutput, &local.Local{})
}

func (a *App) RegisterModule(moduleType ModuleType, module core.Module) error {
	globalName := fmt.Sprintf("%s_%s", moduleType.String(), module.GetName())
	a.registeredModules[globalName] = module
	log.Printf("register module: %s", globalName)
	return nil
}
