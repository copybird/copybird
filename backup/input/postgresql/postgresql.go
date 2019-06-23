package postgres

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"io"
	//"strconv"
	//"strings"
	"text/template"
	//"time"

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

// GetName returns name of module
func (d *PostgresDumper) GetName() string {
	return moduleName
}

// GetConfig returns Config of module
func (d *PostgresDumper) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (d *PostgresDumper) InitPipe(w io.Writer, r io.Reader, cfg interface{}) error {
	d.reader = r
	d.writer = w
	return nil
}

// InitModule initializes module
func (d *PostgresDumper) InitModule(cfg interface{}) error {
	d.config = cfg.(*Config)
	conn, err := sql.Open("postgres", DSN)
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
	/*
		if err := d.dumpDatabase(); err != nil {
			return err
		}
		if err := d.template.Execute(d.writer, d.data); err != nil {
			return err
		}
	*/
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

/*
func (d *PostgresDumper) getTableSchema(tableName string) ([]tableScheme, []sequenceScheme, error) {

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

func (d *PostgresDumper) tableSequenceDump(tableName string, schemas []sequenceScheme) string {
	var sequence []string
	for _, schema := range schemas {
		if schema.name != "" {
			sequence = append(sequence, fmt.Sprintf("drop sequence IF EXISTS %s;\ncreate sequence %s;", schema.name, schema.name))
		}
	}
	return fmt.Sprintf(strings.Join(sequence, ";"))
}

func (d *PostgresDumper) tableSchemeDump(tableName string, schemas []tableScheme) string {

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

func (d *PostgresDumper) getTableData(name string) (string, error) {

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
func (d *PostgresDumper) getServerVersion() (string, error) {
	var version sql.NullString
	if err := d.conn.QueryRow("SELECT version()").Scan(&version); err != nil {
		return version.String, nil
	}
	return version.String, nil
}
*/
