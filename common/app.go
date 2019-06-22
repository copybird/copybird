package common

import (
	"flag"
	"github.com/copybird/copybird/core"
	//"log"

	//"github.com/spf13/cobra"
)

type App struct {
	registeredModules map[string]core.Module
}

func NewApp() *App {
	return &App{
		registeredModules: make(map[string]core.Module),
	}
}

func (a *App) Run() error {
	a.Setup()
	flag.Parse()
	//var cmdBackup = &cobra.Command{
	//	Use:   "backup",
	//	Short: "Start new backup",
	//	Long:  ``,
	//	Args:  cobra.MinimumNArgs(1),
	//	Run: func(cmd *cobra.Command, args []string) {
	//		log.Printf("cmd %#v %s", cmd, args)
	//		a.DoBackup()
	//	},
	//}
	//
	//var rootCmd = &cobra.Command{Use: "copybird"}
	//rootCmd.AddCommand(cmdBackup)
	//return rootCmd.Execute()
	return nil
}

func (a *App) DoBackup() {

}
