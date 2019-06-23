package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func (a *App) Setup() error {
	a.addFlagString(a.cmdBackup, "config", "")
	a.addFlagString(a.cmdBackup, "connect", "")
	a.addFlagString(a.cmdBackup, "input", "")
	a.addFlagString(a.cmdBackup, "compress", "")
	a.addFlagString(a.cmdBackup, "encrypt", "")
	a.addFlagStrings(a.cmdBackup, "output")
	a.addFlagStrings(a.cmdBackup, "notifier")

	a.addFlagString(a.cmdRestore, "config", "")
	a.addFlagString(a.cmdRestore, "connect", "")
	a.addFlagString(a.cmdRestore, "input", "")
	a.addFlagString(a.cmdRestore, "decompress", "")
	a.addFlagString(a.cmdRestore, "decrypt", "")
	a.addFlagStrings(a.cmdRestore, "output")
	a.addFlagStrings(a.cmdRestore, "notifier")

	return nil
}

func (a *App) addFlagString(cmd *cobra.Command, name string, defaultValue string) {
	a.vars[name] = cmd.Flags().String(name, defaultValue, fmt.Sprintf("env %s", strings.ToUpper(name)))
}

func (a *App) addFlagStrings(cmd *cobra.Command, name string) {
	a.vars[name] = cmd.Flags().StringArray(name, nil, fmt.Sprintf("env %s", strings.ToUpper(name)))
}