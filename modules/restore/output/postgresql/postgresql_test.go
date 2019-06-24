package postgres

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

var rs RestoreOutputPostgresql
var conf = Config{
	DSN: "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=test sslmode=disable",
}

func TestRestoreOutputPostgresql_Run(t *testing.T) {
	err := rs.InitModule(&conf)
	assert.Equal(t, err, nil)

	// TODO: Need parse file, but after implement sql formatter
	//f, err := os.Open("../../../../samples/postgres.sql")
	//assert.Equal(t, err, nil)
	//rs.reader = bufio.NewReader(f)

	f, _ := os.Open("../../../../samples/postgres.sql")
	rs.InitPipe(nil, f)
	err = rs.Run()
	assert.Equal(t, err, nil)
}
