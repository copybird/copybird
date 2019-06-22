package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const authorsSchema = `CREATE TABLE authors (id integer default nextval('authors_id_seq'::regclass) not null constraint authors_pk PRIMARY KEY,first_name character varying(100)  not null ,last_name character varying(100)  not null ,email character varying(100)   ,created timestamp without time zone default now() not null );`
const sequenceSchema = `drop sequence IF EXISTS public.authors_id_seq;
create sequence public.authors_id_seq;`
const authorsData = "(1,'test','test','te@asd.ru','2019-06-22T13:26:37.078767Z'),(2,'vanya','ivanov',NULL,'2019-06-22T14:45:54.81458Z')"

func TestGetTables(t *testing.T) {
	d := &PostgresDumper{}
	c := d.GetConfig().(*Config)
	c.DSN = "host=127.0.0.1 port=5432 user=kbereza password=1qazXSW@ dbname=copybird sslmode=disable"

	require.NoError(t, d.InitModule(c))

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
	assert.Equal(t, authorsData, d.data.Tables[0].Data)
	assert.Equal(t, authorsSchema, d.data.Tables[0].Schema)

	assert.NoError(t, d.Run())
}
