package input

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDumper is struct storing inner properties for mysql backups
type MySQLDumper struct {
	conn *sql.DB
}

// New inilializes new MySQLDumper instance
func New(dsn string) (*MySQLDumper, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return &MySQLDumper{conn: conn}, nil
}
