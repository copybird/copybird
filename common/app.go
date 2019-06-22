package common

import (
	"github.com/copybird/copybird/core"
	"github.com/spf13/cobra"
	"log"

	//"log"

	//"github.com/spf13/cobra"
)

type App struct {
	registeredModules map[string]core.Module
	cmmRoot           *cobra.Command
	cmdBackup         *cobra.Command
}

func NewApp() *App {
	return &App{
		registeredModules: make(map[string]core.Module),
	}
}

func (a *App) Run() error {

	var rootCmd = &cobra.Command{Use: "copybird"}
	a.cmdBackup = &cobra.Command{
		Use:   "backup",
		Short: "Start new backup",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("cmd %#v %s", cmd, args)
			a.DoBackup()
		},
	}
	rootCmd.AddCommand(a.cmdBackup)
	a.Setup()
	return rootCmd.Execute()
}

func (a *App) DoBackup() {

}
