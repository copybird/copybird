package postgres

import (
	"bufio"
	"gotest.tools/assert"
	"os"
	"testing"
)

var rs RestoreOutputPostgresql
var conf = Config{
	DSN: "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=test sslmode=disable",
}

func TestRestoreOutputPostgresql_Run(t *testing.T) {
	err := rs.InitModule(&conf)
	assert.Equal(t, err, nil)

	// TODO: Need parse file, but after implement sql formatter
	f, err := os.Open("../../../../samples/postgres.sql")
	assert.Equal(t, err, nil)
	rs.reader = bufio.NewReader(f)

	//rs.reader = bytes.NewReader([]byte("connect test \n "))
	err = rs.Run()
	if err != nil {
		print(err.Error())
	}
	assert.Equal(t, err, nil)
}
