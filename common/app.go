package common

import (
	"github.com/copybird/copybird/core"
	"github.com/copybird/copybird/modules/backup/compress/gzip"
	"github.com/copybird/copybird/modules/backup/compress/lz4"
	"github.com/copybird/copybird/modules/backup/encrypt/aesgcm"
	"github.com/copybird/copybird/modules/backup/input/mongodb"
	"github.com/copybird/copybird/modules/backup/input/mysql"
	"github.com/copybird/copybird/modules/backup/input/mysqldump"
	postgres "github.com/copybird/copybird/modules/backup/input/postgresql"
	"github.com/copybird/copybird/modules/backup/output/gcp"
	"github.com/copybird/copybird/modules/backup/output/http"
	"github.com/copybird/copybird/modules/backup/output/local"
	"github.com/copybird/copybird/modules/backup/output/s3"
	"github.com/copybird/copybird/modules/backup/output/scp"
	"github.com/spf13/cobra"
	"log"
	//"log"
	//"github.com/spf13/cobra"
)

type App struct {
	cmmRoot    *cobra.Command
	cmdBackup  *cobra.Command
	cmdRestore *cobra.Command
	vars       map[string]interface{}
}

func NewApp() *App {
	return &App{
		vars: make(map[string]interface{}),
	}
}

func (a *App) Run() error {
	a.registerModules()
	var rootCmd = &cobra.Command{Use: "copybird"}
	a.cmdBackup = &cobra.Command{
		Use:   "backup",
		Short: "Start new backup",
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdCallback(a.DoBackup),
	}
	a.cmdRestore = &cobra.Command{
		Use:   "restore",
		Short: "Start new restore",
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdCallback(a.DoRestore),
	}
	rootCmd.AddCommand(a.cmdBackup)
	rootCmd.AddCommand(a.cmdRestore)
	a.Setup()
	return rootCmd.Execute()
}

func cmdCallback(f func() error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := f()
		if err != nil {
			log.Printf("cmd err: %s", err)
		}
	}
}

func (a *App) registerModules() {
	core.RegisterModule(&mysql.BackupInputMysql{})
	core.RegisterModule(&mysqldump.BackupInputMysqlDump{})
	core.RegisterModule(&postgres.BackupInputPostgresql{})
	core.RegisterModule(&mongodb.BackupInputMongodb{})
	core.RegisterModule(&gzip.BackupCompressGzip{})
	core.RegisterModule(&lz4.BackupCompressLz4{})
	core.RegisterModule(&aesgcm.BackupEncryptAesgcm{})
	core.RegisterModule(&gcp.BackupOutputGcp{})
	core.RegisterModule(&http.BackupOutputHttp{})
	core.RegisterModule(&local.BackupOutputLocal{})
	core.RegisterModule(&s3.BackupOutputS3{})
	core.RegisterModule(&scp.BackupOutputScp{})
}

func (a *App) Setup() error {

	a.addFlagString(a.cmdBackup, "config", "f", "", "")
	a.addFlagString(a.cmdBackup, "connect", "c", "", "")
	a.addFlagString(a.cmdBackup, "input", "i", "", "(required)")
	a.addFlagString(a.cmdBackup, "compress", "z", "", "")
	a.addFlagString(a.cmdBackup, "encrypt", "e", "", "")
	a.addFlagString(a.cmdBackup, "output", "o", "", "(required)")
	a.addFlagStrings(a.cmdBackup, "notifier", "n", "")

	a.addFlagString(a.cmdRestore, "config", "f", "", "")
	a.addFlagString(a.cmdRestore, "connect", "c", "", "")
	a.addFlagString(a.cmdRestore, "input", "i", "", "(required)")
	a.addFlagString(a.cmdRestore, "decompress", "z", "", "")
	a.addFlagString(a.cmdRestore, "decrypt", "e", "", "")
	a.addFlagString(a.cmdRestore, "output", "o", "", "(required)")
	a.addFlagStrings(a.cmdRestore, "notifier", "n", "")

	return nil
}

func (a *App) addFlagString(cmd *cobra.Command, name, shortName, defaultValue, comment string) {
	a.vars[name] = cmd.Flags().StringP(name, shortName, defaultValue, comment)
}

func (a *App) addFlagStrings(cmd *cobra.Command, name, shortName, comment string) {
	a.vars[name] = cmd.Flags().StringArrayP(name, shortName, nil, comment)
}
