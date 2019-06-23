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

	_ "github.com/go-sql-driver/mysql"
)

// Module Constants
const GroupName = "backup"
const TypeName = "input"
const ModuleName = "mysql"

type (
	// BackupInputMysql is struct storing inner properties for mysql backups
	BackupInputMysql struct {
		core.Module
		conn           *sql.Tx
		headerTemplate *template.Template
		footerTemplate *template.Template
		tableTemplate  *template.Template
		config         *MySQLConfig
		reader         io.Reader
		writer         io.Writer
	}
	dumpHeader struct {
		Version string
	}
	dumpFooter struct {
		EndTime string
	}
	table struct {
		Name   string
		Schema string
		Data   string
	}
)

// GetGroup returns group
func (m *BackupInputMysql) GetGroup() core.ModuleGroup {
	return GroupName
}

// GetType returns type
func (m *BackupInputMysql) GetType() core.ModuleType {
	return TypeName
}

// GetName returns name of module
func (m *BackupInputMysql) GetName() string {
	return ModuleName
}

// GetConfig returns config of module
func (m *BackupInputMysql) GetConfig() interface{} {
	return &MySQLConfig{
		DSN: "root:root@tcp(localhost:3306)/test",
	}
}

// InitPipe initializes pipe
func (m *BackupInputMysql) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

// InitModule initializes module
func (m *BackupInputMysql) InitModule(cfg interface{}) error {
	m.config = cfg.(*MySQLConfig)
	conn, err := sql.Open("mysql", m.config.DSN)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		return err
	}
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	m.conn = tx
	t, err := template.New("headerTemplate").Parse(headerTemplate)
	if err != nil {
		return err
	}
	m.headerTemplate = t
	t, err = template.New("footerTemplate").Parse(footerTemplate)
	if err != nil {
		return err
	}
	m.footerTemplate = t
	t, err = template.New("tableTemplate").Parse(tableTemplate)
	if err != nil {
		return err
	}
	m.tableTemplate = t
	return nil
}

// Run dumps database
func (m *BackupInputMysql) Run() error {
	return m.dumpDatabase()
}

// Close closes ...
func (m *BackupInputMysql) Close() error {
	return nil
}
func (m *BackupInputMysql) getTables() ([]string, error) {
	tables := []string{}
	rows, err := m.conn.Query("SHOW TABLES")
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
func (m *BackupInputMysql) getTableSchema(name string) (string, error) {
	q := fmt.Sprintf("SHOW CREATE TABLE %s", name)
	var returnTable sql.NullString
	var sqlTable sql.NullString
	if err := m.conn.QueryRow(q).Scan(&returnTable, &sqlTable); err != nil {
		return "", err
	}
	if returnTable.String != name {
		return "", errors.New("wrong table returned")
	}
	return sqlTable.String, nil
}

func (m *BackupInputMysql) getTableData(name string) (string, error) {
	q := fmt.Sprintf("SELECT * FROM %s", name)
	rows, err := m.conn.Query(q)
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
func (m *BackupInputMysql) dumpDatabase() error {
	version, err := m.getServerVersion()
	if err != nil {
		return err
	}
	if err := m.headerTemplate.Execute(m.writer, dumpHeader{Version: version}); err != nil {
		return err
	}
	tables, err := m.getTables()
	if err != nil {
		return err
	}
	for _, tableName := range tables {
		var table table
		table.Name = tableName
		schema, err := m.getTableSchema(tableName)
		if err != nil {
			return err
		}
		table.Schema = schema
		data, err := m.getTableData(tableName)
		if err != nil {
			return err
		}
		table.Data = data
		if err := m.tableTemplate.Execute(m.writer, table); err != nil {
			return err
		}
	}
	return m.footerTemplate.Execute(m.writer, dumpFooter{EndTime: time.Now().String()})
}
func (m *BackupInputMysql) getServerVersion() (string, error) {
	var version sql.NullString
	if err := m.conn.QueryRow("SELECT version()").Scan(&version); err != nil {
		return version.String, nil
	}
	return version.String, nil
}
