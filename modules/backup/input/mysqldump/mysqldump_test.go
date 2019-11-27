package mysqldump

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMysqlDumpRun(t *testing.T) {
	m := &BackupInputMysqlDump{}
	c := m.GetConfig().(*MySQLDumpConfig)
	c.Host = "127.0.0.1"
	c.Port = "3306"
	c.Username = "root"
	c.Password = "root"
	c.Database = "test"

	require.NoError(t, m.InitModule(c))
	buf := bytes.Buffer{}
	assert.NoError(t, m.InitPipe(&buf, nil))
	assert.NoError(t, m.Run(context.TODO()))
	t.Log(buf.String())
}
