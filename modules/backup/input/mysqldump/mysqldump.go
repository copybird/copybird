package mysqldump

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/copybird/copybird/core"

	_ "github.com/go-sql-driver/mysql"
)

// Module Constants
const GroupName = "backup"
const TypeName = "input"
const ModuleName = "mysqldump"

type (
	// BackupInputMysqlDump is struct storing inner properties for mysql backups
	BackupInputMysqlDump struct {
		core.Module
		command string
		args    []string
		config  *MySQLDumpConfig
		reader  io.Reader
		writer  io.Writer
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
func (m *BackupInputMysqlDump) GetGroup() core.ModuleGroup {
	return GroupName
}

// GetType returns type
func (m *BackupInputMysqlDump) GetType() core.ModuleType {
	return TypeName
}

// GetName returns name of module
func (m *BackupInputMysqlDump) GetName() string {
	return ModuleName
}

// GetConfig returns config of module
func (m *BackupInputMysqlDump) GetConfig() interface{} {
	return &MySQLDumpConfig{
		Host:              "127.0.0.1",
		Port:              "3306",
		Username:          "root",
		Password:          "root",
		Database:          "test",
		Routines:          true,
		Events:            true,
		Triggers:          true,
		SingleTransaction: true,
		ColumnStatistics:  false,
	}
}

// InitPipe initializes pipe
func (m *BackupInputMysqlDump) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

// InitModule initializes module
func (m *BackupInputMysqlDump) InitModule(cfg interface{}) error {
	m.config = cfg.(*MySQLDumpConfig)

	command, err := exec.LookPath("mysqldump")
	if err != nil {
		return fmt.Errorf("%s cannot be found", command)
	}
	args := []string{
		fmt.Sprintf("-h%s", m.config.Host),
		fmt.Sprintf("-P%s", m.config.Port),
		fmt.Sprintf("-u%s", m.config.Username),
		fmt.Sprintf("-p%s", m.config.Password),
		fmt.Sprintf("--triggers=%t", m.config.Triggers),
		fmt.Sprintf("--routines=%t", m.config.Routines),
		fmt.Sprintf("--events=%t", m.config.Events),
		fmt.Sprintf("--single-transaction=%t", m.config.SingleTransaction),
		// fmt.Sprintf("--column-statistics=%t", m.config.ColumnStatistics),
		m.config.Database,
	}

	m.command = command
	m.args = args
	return nil
}

// Run dumps database
func (m *BackupInputMysqlDump) Run(ctx context.Context) error {
	cmd := exec.Command(m.command, m.args...)
	cmd.Stdout = m.writer
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
