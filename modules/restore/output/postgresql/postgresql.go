package postgres

import (
	"bufio"
	"database/sql"
	"io"
	"log"
	"strings"

	"github.com/copybird/copybird/core"
	_ "github.com/lib/pq"
)

// Module Constants
const GROUP_NAME = "restore"
const TYPE_NAME = "output"
const MODULE_NAME = "postgresql"

type (
	// BackupInputPostgresql is struct storing inner properties for mysql backups
	RestoreOutputPostgresql struct {
		core.Module
		reader io.Reader
		writer io.Writer
		conn   *sql.DB
		config *Config
	}
)

// GetGroup returns module group
func (r *RestoreOutputPostgresql) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

// GetType returns module type
func (r *RestoreOutputPostgresql) GetType() core.ModuleType {
	return TYPE_NAME
}

// GetName returns name of module
func (r *RestoreOutputPostgresql) GetName() string {
	return MODULE_NAME
}

// GetConfig returns Config of module
func (r *RestoreOutputPostgresql) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initialize reader and writer
func (r *RestoreOutputPostgresql) InitPipe(w io.Writer, rd io.Reader) error {
	r.reader = rd
	r.writer = w
	return nil
}

// InitModule initialize connection to DB
func (r *RestoreOutputPostgresql) InitModule(cfg interface{}) error {
	r.config = cfg.(*Config)
	conn, err := sql.Open("postgres", r.config.DSN)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		return err
	}
	r.conn = conn
	return nil
}

// Run dumps database
func (r *RestoreOutputPostgresql) Run() error {
	reader := bufio.NewReader(r.reader)

	var lines []string
	// TODO: Need validate for SQL-like string here.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return err
		}

		// Comment protection
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
		if len(lines) > 1 {
			line = strings.Join(lines, " ")
		}
		err = r.execute(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close connection to DB.
func (r *RestoreOutputPostgresql) Close() error {
	return r.conn.Close()
}

func (r *RestoreOutputPostgresql) execute(line string) error  {
	// Start transaction
	tx, err := r.conn.Begin()
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
