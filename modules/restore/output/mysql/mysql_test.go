package mysql

import (
	"bufio"
	"gotest.tools/assert"
	"os"
	"testing"
)

var rs RestoreOutputMysql
var conf = MySQLConfig{
	DSN: "root:root@tcp(localhost:3306)/test",
}

func TestRestoreOutputMysql_Run(t *testing.T) {
	err := rs.InitModule(&conf)
	assert.Equal(t, err, nil)

	// TODO: Need parse file, but after implement sql formatter
	f, _ := os.Open("../../../..//samples/mysql.sql")
	rs.reader = bufio.NewReader(f)

	err = rs.Run()
	if err != nil {
		print(err.Error())
	}
	assert.Equal(t, err, nil)
}
