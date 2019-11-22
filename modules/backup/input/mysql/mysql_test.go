package mysql

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const authorsSchema = "CREATE TABLE `authors` (\n" +
	"  `id` int(11) NOT NULL AUTO_INCREMENT,\n" +
	"  `first_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `last_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `email` varchar(100) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `birthdate` date NOT NULL,\n" +
	"  `added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
	"  PRIMARY KEY (`id`),\n" +
	"  UNIQUE KEY `email` (`email`)\n" +
	") ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci"

func TestGetTables(t *testing.T) {
	m := &BackupInputMysql{}
	c := m.GetConfig().(*MySQLConfig)
	c.DSN = "root:root@tcp(localhost:3306)/test"

	require.NoError(t, m.InitModule(c))
	f, err := os.Create("dump.sql")
	assert.NoError(t, err)
	assert.NoError(t, m.InitPipe(f, nil))
	tables, err := m.getTables()
	require.NoError(t, err)
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "authors", tables[0])
	tableSchema, err := m.getTableSchema(tables[0])
	assert.NoError(t, err)
	assert.Equal(t, authorsSchema, tableSchema)
	err = m.Run(context.TODO())
	assert.NoError(t, err)
	assert.NoError(t, os.Remove("dump.sql"))
}

func TestMysqlDump(t *testing.T) {
	m := &BackupInputMysql{}

	c := m.GetConfig().(*MySQLConfig)
	c.DSN = "root:root@tcp(localhost:3306)/test"

	require.NoError(t, m.InitModule(c))
	buf := bytes.Buffer{}
	require.NoError(t, m.InitPipe(&buf, nil))
	assert.NoError(t, m.Run(context.TODO()))
	t.Log(buf.String())
}
