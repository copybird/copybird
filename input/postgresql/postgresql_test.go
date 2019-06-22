package postgres

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const authorsSchema = `CREATE TABLE authors(id integer default nextval('authors_id_seq'::regclass) not null constraint authors_pk PRIMARY KEY,first_name character varying(100)  not null ,last_name character varying(100)  not null ,email character varying(100)   ,created timestamp without time zone default now() not null );
create sequence public.authors_id_seq;`
const authorsData = "(1,'test','test','te@asd.ru','2019-06-22T13:26:37.078767Z')"

func TestGetTables(t *testing.T) {
	d := &PostgresDumper{}
	c := d.GetConfig().(*config)
	c.DSN = "host=127.0.0.1 port=5432 user=postgres dbname=blog sslmode=disable"

	var d2 = []byte("")
	r := bytes.NewBuffer(d2)

	fmt.Println(d.InitPipe(r, nil, c))

	require.NoError(t, d.InitModule(c))

	tables, err := d.getTables()
	assert.NoError(t, err)
	assert.Equal(t, "authors", tables[0])
	assert.Equal(t, "posts", tables[1])

	tableSchema, err := d.getTableSchema(tables[0])
	fmt.Println("err", err, tableSchema)
	assert.NoError(t, err)
	assert.Equal(t, authorsSchema, tableSchema)

	data, err := d.getTableData(tables[0])
	assert.NoError(t, err)
	assert.Equal(t, authorsData, data)

	assert.NoError(t, d.dumpDatabase())
	assert.Equal(t, authorsData, d.data.Tables[0].Data)
	assert.Equal(t, authorsSchema, d.data.Tables[0].Schema)

	assert.NoError(t, d.Run())

	err = ioutil.WriteFile(`E:\OpenServer\domains\projects\go\src\github.com\copybird\copybird\test.sql`, r.Bytes(), 0644)
	assert.NoError(t, err)
}
