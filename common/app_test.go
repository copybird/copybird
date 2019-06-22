package common

import (
	"io"
	"log"
	"sync"
	"testing"

	"github.com/copybird/copybird/compress/gzip"
	"github.com/copybird/copybird/encryption/aesgcm"
	"github.com/copybird/copybird/input/mysql"
	"github.com/copybird/copybird/output/local"

	"gotest.tools/assert"
)

func TestAppRun(t *testing.T) {
	modMysql := new(mysql.MySQLDumper)
	modMysqlCfg := modMysql.GetConfig()
	assert.Assert(t, modMysqlCfg != nil)
	assert.Assert(t, modMysql.GetConfig() != modMysql.GetConfig())
	assert.NilError(t, modMysql.InitModule(modMysqlCfg))

	modGzip := new(gzip.Compress)
	modGzipCfg := modGzip.GetConfig()
	assert.Assert(t, modGzipCfg != nil)
	assert.Assert(t, modGzip.GetConfig() != modGzip.GetConfig())
	assert.NilError(t, modMysql.InitModule(modGzipCfg))

	modAesgcm := new(aesgcm.EncryptionAESGCM)
	modAesgcmCfg := modAesgcm.GetConfig()
	assert.Assert(t, modAesgcm != nil)
	assert.Assert(t, modAesgcm.GetConfig() != modAesgcm.GetConfig())
	assert.NilError(t, modAesgcm.InitModule(modAesgcmCfg))

	modOutput := new(local.Local)
	modOutputCfg := modOutput.GetConfig()
	assert.Assert(t, modOutput != nil)
	assert.Assert(t, modOutput.GetConfig() != modOutput.GetConfig())
	assert.NilError(t, modOutput.InitModule(modOutputCfg))

	inputReader, inputWriter := io.Pipe()
	gzipReader, gzipWriter := io.Pipe()
	aesgcmReader, aesgcmWriter := io.Pipe()

	assert.NilError(t, modMysql.InitPipe(inputWriter, nil))
	assert.NilError(t, modGzip.InitPipe(gzipWriter, inputReader))
	assert.NilError(t, modAesgcm.InitPipe(aesgcmWriter, gzipReader))
	assert.NilError(t, modOutput.InitPipe(nil, aesgcmReader))

	wg := sync.WaitGroup{}
	wg.Add(4)

	chErr := make(chan error, 10)

	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := modMysql.Run(); err != nil {
			log.Printf("%s err: %s", modMysql.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := modGzip.Run(); err != nil {
			log.Printf("%s err: %s", modGzip.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := modAesgcm.Run(); err != nil {
			log.Printf("%s err: %s", modAesgcm.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := modOutput.Run(); err != nil {
			log.Printf("%s err: %s", modOutput.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	wg.Wait()

	err, ok := <-chErr
	if ok && err != nil {
		log.Printf("pipe err: %s", err)
	}

}
