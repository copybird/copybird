package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/copybird/copybird/core"
	_ "github.com/lib/pq"
)

const moduleName = "postgresql"

type (
	// PostgresDumper is struct storing inner properties for mysql backups
	PostgresDumper struct {
		core.PipeComponent
		conn     *sql.DB
		data     dbDump
		template *template.Template
		config   *config
		reader   io.Reader
		writer   io.Writer
	}
	dbDump struct {
		Version string
		Tables  []table
		EndTime string
	}
	table struct {
		Name     string
		Schema   string
		Data     string
		DBScheme string
	}
)

// GetName returns name of module
func (d *PostgresDumper) GetName() string {
	return moduleName
}

// GetConfig returns config of module
func (d *PostgresDumper) GetConfig() interface{} {
	return &config{}
}

// InitPipe initializes pipe
func (d *PostgresDumper) InitPipe(w io.Writer, r io.Reader, cfg interface{}) error {
	d.reader = r
	d.writer = w
	return nil
}

// InitModule initializes module
func (d *PostgresDumper) InitModule(cfg interface{}) error {
	d.config = cfg.(*config)
	conn, err := sql.Open("postgres", d.config.DSN)
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
func (d *PostgresDumper) Run() error {
	if err := d.dumpDatabase(); err != nil {
		return err
	}
	if err := d.template.Execute(d.writer, d.data); err != nil {
		return err
	}

	return nil
}

// Close closes ...
func (d *PostgresDumper) Close() error {
	return nil
}
func (d *PostgresDumper) getTables() ([]string, error) {
	var (
		tables    []string
		tableType = "BASE TABLE"
	)
	rows, err := d.conn.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%s' AND table_type='%s'", defaultSchemaName, tableType))
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

func (d *PostgresDumper) getTableSchema(tableName string) (string, error) {

	var (
		columns     []tableColumn
		sequence    []string
		tableSchema = fmt.Sprintf("CREATE TABLE %s", tableName)
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
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName, columnDefault, dataType, characterMaximumLength, isNullable, constraintName, constraintType, sequence sql.NullString
		if err := rows.Scan(&tableName, &columnName, &columnDefault, &dataType, &characterMaximumLength, &isNullable, &constraintName, &constraintType, &sequence); err != nil {
			return tableSchema, err
		}
		columns = append(columns, tableColumn{
			columnName:             columnName.String,
			columnDefault:          columnDefault.String,
			dataType:               dataType.String,
			characterMaximumLength: characterMaximumLength.String,
			isNullable:             isNullable.String,
			constraintName:         constraintName.String,
			constraintType:         constraintType.String,
			sequence:               sequence.String,
		})
	}

	var tableColumns []string
	for _, column := range columns {
		var (
			defaultValue, isNull, constraint string
			columnType                       = column.dataType
		)
		if column.columnDefault != "" {
			defaultValue = fmt.Sprintf("default %s", column.columnDefault)
		}

		if column.characterMaximumLength != "" {
			columnType = fmt.Sprintf("%s(%s)", column.dataType, column.characterMaximumLength)
		}

		if column.isNullable == "NO" {
			isNull = "not null"
		}

		if column.constraintName != "" {
			constraint = fmt.Sprintf("constraint %s %s", column.constraintName, column.constraintType)
		}

		tableColumns = append(tableColumns, fmt.Sprintf("%s %s %s %s %s", column.columnName, columnType, defaultValue, isNull, constraint))
		if column.sequence != "" {
			sequence = append(sequence, fmt.Sprintf("create sequence %s", column.sequence))
		}
	}

	tableSchema += fmt.Sprintf("(%s);\n%s;", strings.Join(tableColumns, ","), strings.Join(sequence, ";"))
	return tableSchema, nil
}

func (d *PostgresDumper) getTableData(name string) (string, error) {
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
func (d *PostgresDumper) dumpDatabase() error {
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
		table.DBScheme = defaultSchemaName
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
	t, err := template.New("postgresqlbackup").Parse(dumpTemplate)
	if err != nil {
		return err
	}
	d.template = t
	d.data = dump
	return nil

}
func (d *PostgresDumper) getServerVersion() (string, error) {
	var version sql.NullString
	if err := d.conn.QueryRow("SELECT version()").Scan(&version); err != nil {
		return version.String, nil
	}
	return version.String, nil
}
