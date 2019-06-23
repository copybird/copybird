package postgres

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"fmt"
	"io"

	"text/template"

	"github.com/copybird/copybird/core"
	_ "github.com/lib/pq"
)

// Module Constants
const GROUP_NAME = "backup"
const TYPE_NAME = "input"
const moduleName = "postgresql"

type (
	// BackupInputPostgresql is struct storing inner properties for mysql backups
	BackupInputPostgresql struct {
		core.Module
		conn     *sql.DB
		data     dbDump
		template *template.Template
		config   *Config
		reader   io.Reader
		writer   io.Writer
	}
	dbDump struct {
		Version  string
		Tables   []table
		EndTime  string
		DBScheme string
	}
	table struct {
		Name           string
		Schema         string
		SequenceScheme string
		Data           string
		DBScheme       string
	}
)

// GetGroup returns group
func (m *BackupInputPostgresql) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

// GetType returns type
func (m *BackupInputPostgresql) GetType() core.ModuleType {
	return TYPE_NAME
}

// GetName returns name of module
func (m *BackupInputPostgresql) GetName() string {
	return moduleName
}

// GetConfig returns Config of module
func (m *BackupInputPostgresql) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (m *BackupInputPostgresql) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

// InitModule initializes module
func (m *BackupInputPostgresql) InitModule(cfg interface{}) error {
	m.config = cfg.(*Config)
	conn, err := sql.Open("postgres", m.config.DSN)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		return err
	}
	m.conn = conn
	return nil
}

// Run dumps database
func (m *BackupInputPostgresql) Run() error {

	if err := m.dumpDatabase(); err != nil {
		return err
	}
	if err := m.template.Execute(m.writer, m.data); err != nil {
		return err
	}
	return nil
}

// Close closes ...
func (m *BackupInputPostgresql) Close() error {
	return nil
}
func (m *BackupInputPostgresql) getTables() ([]string, error) {
	var (
		tables    []string
		tableType = "BASE TABLE"
	)
	rows, err := m.conn.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%s' AND table_type='%s'", defaultSchemaName, tableType))
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

func (d *BackupInputPostgresql) getTableSchema(tableName string) ([]tableScheme, []sequenceScheme, error) {

	var (
		columns  []tableScheme
		sequence []sequenceScheme
	)
	rows, err := d.conn.Query(`
	select cln.table_name,
		   cln.column_name,
		   cln.column_default,
		   cln.data_type,
		   cln.character_maximum_length,
		   cln.is_nullable,
		   tc.constraint_name,
		   tc.constraint_type,
		   pg_get_serial_sequence($2, cln.column_name) AS sequence
	from INFORMATION_SCHEMA.COLUMNS cln
			 LEFT JOIN INFORMATION_SCHEMA.constraint_column_usage ctu ON ctu.table_schema = cln.table_schema
		AND ctu.column_name = cln.column_name
		and ctu.table_name = cln.table_name
			 LEFT JOIN INFORMATION_SCHEMA.table_constraints tc ON tc.constraint_schema = cln.table_schema
		and tc.table_name = cln.table_name
		and tc.constraint_name = ctu.constraint_name
	where cln.table_schema = $1
	  and cln.table_name = $2`, defaultSchemaName, tableName)
	if err != nil {
		return columns, sequence, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName, columnDefault, dataType, characterMaximumLength, isNullable, constraintName, constraintType, sequenceName sql.NullString
		if err := rows.Scan(&tableName, &columnName, &columnDefault, &dataType, &characterMaximumLength, &isNullable, &constraintName, &constraintType, &sequenceName); err != nil {
			return columns, sequence, err
		}
		columns = append(columns, tableScheme{
			columnName:             columnName.String,
			columnDefault:          columnDefault.String,
			dataType:               dataType.String,
			characterMaximumLength: characterMaximumLength.String,
			isNullable:             isNullable.String,
			constraintName:         constraintName.String,
			constraintType:         constraintType.String,
			sequence:               sequenceName.String,
		})
		sequence = append(sequence, sequenceScheme{name: sequenceName.String})
	}

	return columns, sequence, nil
}

func (d *BackupInputPostgresql) tableSequenceDump(tableName string, schemas []sequenceScheme) string {
	var sequence []string
	for _, schema := range schemas {
		if schema.name != "" {
			sequence = append(sequence, fmt.Sprintf("drop sequence IF EXISTS %s;\ncreate sequence %s;", schema.name, schema.name))
		}
	}
	return fmt.Sprintf(strings.Join(sequence, ";"))
}

func (d *BackupInputPostgresql) tableSchemeDump(tableName string, schemas []tableScheme) string {

	var tableColumns []string
	for _, schema := range schemas {
		var defaultValue, isNull, constraint string
		var columnType = schema.dataType
		if schema.columnDefault != "" {
			defaultValue = fmt.Sprintf("default %s", schema.columnDefault)
		}

		if schema.characterMaximumLength != "" {
			columnType = fmt.Sprintf("%s(%s)", schema.dataType, schema.characterMaximumLength)
		}

		if schema.isNullable == "NO" {
			isNull = "not null"
		}

		if schema.constraintName != "" {
			constraint = fmt.Sprintf("constraint %s %s", schema.constraintName, schema.constraintType)
		}
		tableColumns = append(tableColumns, fmt.Sprintf("%s %s %s %s %s", schema.columnName, columnType, defaultValue, isNull, constraint))
	}

	return fmt.Sprintf("CREATE TABLE %s (%s);", tableName, strings.Join(tableColumns, ","))
}

func (d *BackupInputPostgresql) getTableData(name string) (string, error) {

	rows, err := d.conn.Query(fmt.Sprintf(`SELECT * FROM %s`, name))
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
		var (
			scanData = make([]sql.RawBytes, len(columns))
			pointers = make([]interface{}, len(columns))
		)

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
func (d *BackupInputPostgresql) dumpDatabase() error {
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
	dump.DBScheme = defaultSchemaName
	dump.EndTime = time.Now().String()

	for _, tableName := range tables {
		var table = table{
			Name:     tableName,
			DBScheme: defaultSchemaName,
		}

		tableSchema, sequenceSchema, err := d.getTableSchema(tableName)
		if err != nil {
			return err
		}

		data, err := d.getTableData(tableName)
		if err != nil {
			return err
		}

		table.Schema = d.tableSchemeDump(tableName, tableSchema)
		table.SequenceScheme = d.tableSequenceDump(tableName, sequenceSchema)
		table.Data = data
		dump.Tables = append(dump.Tables, table)
	}

	t, err := template.New("postgresqlbackup").Parse(dumpTemplate)
	if err != nil {
		return err
	}

	d.template = t
	d.data = dump
	return nil

}
func (d *BackupInputPostgresql) getServerVersion() (string, error) {
	var version sql.NullString
	if err := d.conn.QueryRow("SELECT version()").Scan(&version); err != nil {
		return version.String, nil
	}
	return version.String, nil
}
