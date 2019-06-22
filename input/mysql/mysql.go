package input

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

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
		var table sql.NullString
		if err := rows.Scan(&table); err != nil {
			return tables, err
		}
		tables = append(tables, table.String)
	}
	return tables, rows.Err()
}
func (d *MySQLDumper) getTableSchema(name string) (string, error) {
	q := fmt.Sprintf("SHOW CREATE TABLE %s", name)
	var returnTable sql.NullString
	var sqlTable sql.NullString
	if err := d.conn.QueryRow(q).Scan(&returnTable, &sqlTable); err != nil {
		return "", err
	}
	if returnTable.String != name {
		return "", errors.New("wrong table returned")
	}
	return sqlTable.String, nil
}

func (d *MySQLDumper) getTableData(name string) (string, error) {
	q := fmt.Sprintf("SELECT * FROM %s", name)
	rows, err := d.conn.Query(q)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	if len(columns) == 0 {
		return "", fmt.Errorf("no columns in table %s", name)
	}
	var data []string
	for rows.Next() {
		scanData := make([]sql.RawBytes, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range scanData {
			pointers[i] = &scanData[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return "", err
		}
		rowData := make([]string, len(columns))
		for i, v := range scanData {
			if v != nil {
				if _, err := strconv.Atoi(string(v)); err == nil {
					rowData[i] = string(v)
				} else {
					rowData[i] = fmt.Sprintf("'%s'", strings.Replace(string(v), "'", "\\'", -1))
				}
			} else {
				rowData[i] = "NULL"
			}
			json.Unmarshal(v, pointers)
		}
		data = append(data, fmt.Sprintf("(%s)", strings.Join(rowData, ",")))

	}
	return strings.Join(data, ","), rows.Err()

}
