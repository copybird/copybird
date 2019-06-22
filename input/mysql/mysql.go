package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/copybird/copybird/core"
	// We need this shit
	_ "github.com/go-sql-driver/mysql"
)

const MODULE_NAME = "mysql"

type (
	// MySQLDumper is struct storing inner properties for mysql backups
	MySQLDumper struct {
		core.PipeComponent
		conn     *sql.DB
		data     dbDump
		template *template.Template
		config   *MySQLConfig
		reader   io.Reader
		writer   io.Writer
	}
	dbDump struct {
		Version string
		Tables  []table
		EndTime string
	}
	table struct {
		Name   string
		Schema string
		Data   string
	}
)

// GetName returns name of module
func (d *MySQLDumper) GetName() string {
	return MODULE_NAME
}

// GetConfig returns config of module
func (d *MySQLDumper) GetConfig() interface{} {
	return &MySQLConfig{}
}

// InitPipe initializes pipe
func (d *MySQLDumper) InitPipe(w io.Writer, r io.Reader, cfg interface{}) error {
	d.reader = r
	d.writer = w
	return nil
}

// InitModule initializes module
func (d *MySQLDumper) InitModule(cfg interface{}) error {
	d.config = cfg.(*MySQLConfig)
	conn, err := sql.Open("mysql", d.config.DSN)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		return err
	}
	d.conn = conn
	return nil
}

// Run dumps database
func (d *MySQLDumper) Run() error {
	if err := d.dumpDatabase(); err != nil {
		return err
	}
	if err := d.template.Execute(d.writer, d.data); err != nil {
		return err
	}

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
func (d *MySQLDumper) dumpDatabase() error {
	var dump dbDump
	version, err := d.getServerVersion()
	if err != nil {
		return err
	}
	tables, err := d.getTables()
	if err != nil {
		return err
	}
	dump.Version = version
	for _, tableName := range tables {
		var table table
		table.Name = tableName
		schema, err := d.getTableSchema(tableName)
		if err != nil {
			return err
		}
		table.Schema = schema
		data, err := d.getTableData(tableName)
		if err != nil {
			return err
		}
		table.Data = data
		dump.Tables = append(dump.Tables, table)
	}
	dump.EndTime = time.Now().String()
	t, err := template.New("mysqlbackup").Parse(dumpTemplate)
	if err != nil {
		return err
	}
	d.template = t
	d.data = dump
	return nil

}
func (d *MySQLDumper) getServerVersion() (string, error) {
	var version sql.NullString
	if err := d.conn.QueryRow("SELECT version()").Scan(&version); err != nil {
		return version.String, nil
	}
	return version.String, nil
}
