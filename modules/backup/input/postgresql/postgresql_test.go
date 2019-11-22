package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const authorsSchema = `CREATE TABLE authors (id integer default nextval('authors_id_seq'::regclass) not null constraint authors_pk PRIMARY KEY,first_name character varying(100)  not null ,last_name character varying(100)  not null ,email character varying(100)   ,created timestamp without time zone default now() not null );`
const sequenceSchema = `drop sequence IF EXISTS public.authors_id_seq;
create sequence public.authors_id_seq;`
const authorsData = "(1,'test','test','te@asd.ru','2019-06-22T13:26:37.078767Z'),(2,'vanya','ivanov',NULL,'2019-06-22T14:45:54.81458Z')"

func TestGetTables(t *testing.T) {
	d := &BackupInputPostgresql{}
	c := d.GetConfig().(*Config)
	c.DSN = "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=test sslmode=disable"

	f, err := os.Create("dump.sql")
	assert.NoError(t, err)
	require.NoError(t, d.InitModule(c))
	assert.NoError(t, d.InitPipe(f, nil))

	tables, err := d.getTables()
	assert.NoError(t, err)
	assert.Equal(t, "authors", tables[0])
	assert.Equal(t, "posts", tables[1])

	tableSchema, seqSchema, err := d.getTableSchema(tables[0])
	fmt.Println("err", err, tableSchema, seqSchema)

	assert.NoError(t, err)
	assert.Equal(t, authorsSchema, d.tableSchemeDump(tables[0], tableSchema))
	assert.Equal(t, sequenceSchema, d.tableSequenceDump(tables[0], seqSchema))

	data, err := d.getTableData(tables[0])
	assert.NoError(t, err)
	assert.Equal(t, authorsData, data)

	assert.NoError(t, d.dumpDatabase())

	assert.NoError(t, d.Run(context.TODO()))

	assert.NoError(t, os.Remove("dump.sql"))
}
