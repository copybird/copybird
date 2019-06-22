package common

import (
	"fmt"
	compress_gzip "github.com/copybird/copybird/compress/gzip"
	"github.com/copybird/copybird/core"
	"github.com/copybird/copybird/encryption/aesgcm"
	"github.com/copybird/copybird/input/mysql"
	"github.com/copybird/copybird/output/local"
	"github.com/copybird/copybird/output/s3"

	//"log"
	"reflect"
	//"strings"

	"github.com/iancoleman/strcase"

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
	a.RegisterModule(ModuleTypeOutput, &s3.S3{})
}

func (a *App) RegisterModule(moduleType ModuleType, module core.Module) error {
	moduleGlobalName := fmt.Sprintf("%s_%s", moduleType.String(), module.GetName())
	a.registeredModules[moduleGlobalName] = module
	//log.Printf("register module: %s", moduleGlobalName)
	cfg := module.GetConfig()
	cfgValue := reflect.Indirect(reflect.ValueOf(cfg))
	cfgType := cfgValue.Type()
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		name := strcase.ToSnake(field.Name)
		argName := fmt.Sprintf("%s_%s", moduleGlobalName, name)
		//envName := strings.ToUpper(argName)

		//log.Printf("argName: %s, envName: %s", argName, envName)
		switch (field.Type.Kind()) {
		case reflect.Int:
			a.cmdBackup.Flags().Int64(argName, cfgValue.Field(i).Int(), argName)
		case reflect.String:
			a.cmdBackup.Flags().String(argName, cfgValue.Field(i).String(), argName)
		case reflect.Bool:
			a.cmdBackup.Flags().Bool(argName, cfgValue.Field(i).Bool(), argName)
		}
	}
	return nil
}
