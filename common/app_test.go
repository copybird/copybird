package common

import (
	"io"
	"sync"
	"testing"

	"github.com/copybird/copybird/modules/backup/compress/gzip"
	"github.com/copybird/copybird/modules/backup/encrypt/aesgcm"
	"github.com/copybird/copybird/modules/backup/input/mysql"
	"github.com/copybird/copybird/modules/backup/output/local"

	"gotest.tools/assert"
)

func TestAppRun(t *testing.T) {
	modMysql := new(mysql.BackupInputMysql)
	modMysqlCfg := modMysql.GetConfig()
	assert.Assert(t, modMysqlCfg != nil)
	assert.Assert(t, modMysql.GetConfig() != modMysql.GetConfig())
	assert.NilError(t, modMysql.InitModule(modMysqlCfg))

	modGzip := new(gzip.BackupCompressGzip)
	modGzipCfg := modGzip.GetConfig()
	assert.Assert(t, modGzipCfg != nil)
	assert.Assert(t, modGzip.GetConfig() != modGzip.GetConfig())
	assert.NilError(t, modMysql.InitModule(modGzipCfg))

	modAesgcm := new(aesgcm.BackupEncryptAesgcm)
	modAesgcmCfg := modAesgcm.GetConfig()
	assert.Assert(t, modAesgcm != nil)
	assert.Assert(t, modAesgcm.GetConfig() != modAesgcm.GetConfig())
	assert.NilError(t, modAesgcm.InitModule(modAesgcmCfg))

	modOutput := new(local.BackupOutputLocal)
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
			t.Logf("%s err: %s", modMysql.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := modGzip.Run(); err != nil {
			t.Logf("%s err: %s", modGzip.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := modAesgcm.Run(); err != nil {
			t.Logf("%s err: %s", modAesgcm.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := modOutput.Run(); err != nil {
			t.Logf("%s err: %s", modOutput.GetName(), err)
			chErr <- err
		}
	}(&wg, chErr)

	wg.Wait()

	for {
		err, ok := <-chErr
		if !ok {
			break
		}
		if ok && err != nil {
			t.Logf("pipe err: %s", err)
		}
	}

}
