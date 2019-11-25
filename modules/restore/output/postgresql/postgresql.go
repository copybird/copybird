package postgres

import (
	"context"
	"database/sql"
	"io"

	"github.com/copybird/copybird/core"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	"github.com/xwb1989/sqlparser"
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
func (r *RestoreOutputPostgresql) Run(ctx context.Context) error {
	tokenizer := sqlparser.NewTokenizer(r.reader)

	for {
		stmt, err := sqlparser.ParseNext(tokenizer)
		if err != nil {
			spew.Dump(err)
		}
		if err == io.EOF {
			break
		}
		spew.Dump(stmt)
		spew.Dump(sqlparser.String(stmt))
		switch stmt.(type) {
		case *sqlparser.Select, *sqlparser.Insert, *sqlparser.DBDDL, *sqlparser.DDL:
			if _, err := r.conn.Exec(sqlparser.String(stmt)); err != nil {
				return err
			}
		default:
			continue
		}

	}
	return nil
}

// Close connection to DB.
func (r *RestoreOutputPostgresql) Close() error {
	return r.conn.Close()
}

func (r *RestoreOutputPostgresql) execute(line string) error {
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
