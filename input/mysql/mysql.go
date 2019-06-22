package input

import (
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/copybird/copybird/core"
	// We need this shit
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDumper is struct storing inner properties for mysql backups
type MySQLDumper struct {
	core.PipeComponent
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

// Init initializes ...
func (d *MySQLDumper) Init(w io.Writer, r io.Reader) error {
	return nil
}

// Run dumps database
func (d *MySQLDumper) Run() error {
	tables, err := d.getTables()
	if err != nil {
		return err
	}
	_ = tables
	return nil
}

// Close closes ...
func (d *MySQLDumper) Close() error {
	return nil
}
func (d *MySQLDumper) getTables() ([]string, error) {
	tables := []string{}
	rows, err := d.conn.Query("SHOW TABLES")
	if err != nil {
		return tables, err
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return tables, err
		}
		tables = append(tables, table)
	}
	return tables, rows.Err()
}
func (d *MySQLDumper) getTableSchema(name string) (string, error) {
	q := fmt.Sprintf("SHOW CREATE TABLE %s", name)
	var returnTable string
	var sqlTable string
	if err := d.conn.QueryRow(q).Scan(&returnTable, &sqlTable); err != nil {
		return "", err
	}
	if returnTable != name {
		return "", errors.New("wrong table returned")
	}
	return sqlTable, nil
}
