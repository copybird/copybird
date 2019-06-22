package common

import (
	"log"

	"github.com/spf13/cobra"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	var cmdBackup = &cobra.Command{
		Use:   "backup",
		Short: "Start new backup",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("cmd %#v %s", cmd, args)
			a.DoBackup()
		},
	}

	var rootCmd = &cobra.Command{Use: "copybird"}
	rootCmd.AddCommand(cmdBackup)
	return rootCmd.Execute()
}

func (a *App) DoBackup() {

}
