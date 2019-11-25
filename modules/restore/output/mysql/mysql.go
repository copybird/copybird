package mysql

import (
	"context"
	"database/sql"
	"io"

	"github.com/copybird/copybird/core"
	"github.com/davecgh/go-spew/spew"
	"github.com/xwb1989/sqlparser"

	_ "github.com/go-sql-driver/mysql"
)

// Module Constants
const (
	GROUP_NAME  = "restore"
	TYPE_NAME   = "output"
	MODULE_NAME = "mysql"
)

type (
	// RestoreOutputMysql is struct storing inner properties for mysql backups
	RestoreOutputMysql struct {
		core.Module
		conn   *sql.DB
		config *MySQLConfig
		reader io.Reader
		writer io.Writer
	}
)

// GetGroup returns group
func (m *RestoreOutputMysql) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

// GetType returns type
func (m *RestoreOutputMysql) GetType() core.ModuleType {
	return TYPE_NAME
}

// GetName returns name of module
func (m *RestoreOutputMysql) GetName() string {
	return MODULE_NAME
}

// GetConfig returns config of module
func (m *RestoreOutputMysql) GetConfig() interface{} {
	return &MySQLConfig{
		DSN: "root:root@tcp(localhost:3306)/test",
	}
}

// InitPipe initializes pipe
func (m *RestoreOutputMysql) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

// InitModule initializes module
func (m *RestoreOutputMysql) InitModule(cfg interface{}) error {
	m.config = cfg.(*MySQLConfig)
	conn, err := sql.Open("mysql", m.config.DSN)
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
func (m *RestoreOutputMysql) Run(ctx context.Context) error {
	return m.RestoreDatabase()
}

// RestoreDatabase restores db
func (m *RestoreOutputMysql) RestoreDatabase() error {
	tokenizer := sqlparser.NewTokenizer(m.reader)

	for {
		stmt, err := sqlparser.ParseNext(tokenizer)
		if err != nil {
			spew.Dump(err)
		}
		if err == io.EOF {
			break
		}
		switch stmt.(type) {
		case *sqlparser.Select, *sqlparser.Insert, *sqlparser.DBDDL, *sqlparser.DDL:
			if _, err := m.conn.Exec(sqlparser.String(stmt)); err != nil {
				return err
			}
		default:
			continue
		}

	}
	return nil
}

// Close connection to DB.
func (m *RestoreOutputMysql) Close() error {
	return m.conn.Close()
}

func (m *RestoreOutputMysql) execute(line string) error {
	// Start transaction
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}

	// Execute transaction
	_, err = tx.Exec(line)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
