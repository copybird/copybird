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
