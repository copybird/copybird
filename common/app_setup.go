package common

import (
	"reflect"
	"log"

	"github.com/copybird/copybird/core"
	"github.com/copybird/copybird/input/mysql"
	compress_gzip "github.com/copybird/copybird/compress/gzip"
	decompress_gzip "github.com/copybird/copybird/decompress/gzip"
	// lz4_compress "github.com/copybird/copybird/compress/lz4"
	// lz4_decompress "github.com/copybird/copybird/decompress/lz4"
)

func (a *App) Setup() error {
	a.RegisterModule(&mysql.MySQLDumper{})
	a.RegisterModule(&compress_gzip.Compress{})
	a.RegisterModule(&decompress_gzip.Decompress{})
	return nil
}

func (a *App) RegisterModule(module core.Module) error {
	moduleType := reflect.TypeOf(module)
	log.Printf("register module: %s", moduleType.PkgPath())
	return nil
}
