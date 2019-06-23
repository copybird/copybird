package mysql

import (
	"bufio"
	"database/sql"
	"github.com/copybird/copybird/core"
	"io"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Module Constants
const (
	GROUP_NAME = "restore"
	TYPE_NAME = "output"
	MODULE_NAME = "mysql"
)

type (
	// RestoreOutputMysql is struct storing inner properties for mysql backups
	RestoreOutputMysql struct {
		core.Module
		conn           *sql.DB
		config         *MySQLConfig
		reader         io.Reader
		writer         io.Writer
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
func (m *RestoreOutputMysql) Run() error {
	return m.RestoreDatabase()
}

// RestoreDatabase restores db
func (m *RestoreOutputMysql) RestoreDatabase() error {
	reader := bufio.NewReader(m.reader)

	var lines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if strings.HasPrefix(line, "--") {
			continue
		}
		if strings.HasPrefix(line, "/*") {
			continue
		}

		line = strings.TrimSuffix(line, "\n")
		if line == "" {
			continue
		}

		if !strings.HasSuffix(line, ";") {
			lines = append(lines, line)
			continue
		}

		if len(lines) > 0 {
			lines = append(lines, line)
			line = strings.Join(lines, " ")
			lines = []string{}
		}

		err = m.execute(line)
		if err != nil {
			return err
		}

	}
	return nil
}

// Close connection to DB.
func (m *RestoreOutputMysql) Close() error {
	return m.conn.Close()
}

func (m *RestoreOutputMysql) execute(line string) error  {
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
