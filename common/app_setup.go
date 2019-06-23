package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type ModuleType int

func (m ModuleType) String() string {
	return [...]string{
		"connect",
		"input",
		"compress",
		"encrypt",
		"backup_output",

		"notifier",
	}[m]
}

func (a *App) Setup() error {
	a.addFlagString(a.cmdBackup, "config", "")
	a.addFlagString(a.cmdBackup, "input", "mysql")
	a.addFlagString(a.cmdBackup, "compress", "")
	a.addFlagString(a.cmdBackup, "encrypt", "")
	a.addFlagString(a.cmdBackup, "output", "local")
	a.addFlagStrings(a.cmdBackup, "notifier")
	return nil
}

func (a *App) addFlagString(cmd *cobra.Command, name string, defaultValue string) {
	a.vars[name] = cmd.Flags().String(name, defaultValue, fmt.Sprintf("env %s", strings.ToUpper(name)))
}

func (a *App) addFlagStrings(cmd *cobra.Command, name string) {
	a.vars[name] = cmd.Flags().StringArray(name, nil, fmt.Sprintf("env %s", strings.ToUpper(name)))
}
//
//func (a *App) addFlagInt64(cmd *cobra.Command, name string, defaultValue int64) {
//	a.vars[name] = cmd.Flags().Int64(name, defaultValue, fmt.Sprintf("env %s", strings.ToUpper(name)))
//}
//
//func (a *App) addFlagBool(cmd *cobra.Command, name string, defaultValue bool) {
//	a.vars[name] = cmd.Flags().Bool(name, defaultValue, fmt.Sprintf("env %s", strings.ToUpper(name)))
//}

//func (a *App) RegisterModule(moduleType ModuleType, module core.Module) error {
//	moduleGlobalName := fmt.Sprintf("%s_%s", moduleType.String(), module.GetName())
//	a.modulesBackup[moduleGlobalName] = module
//	cfg := module.GetConfig()
//	cfgValue := reflect.Indirect(reflect.ValueOf(cfg))
//	cfgType := cfgValue.Type()
//	for i := 0; i < cfgType.NumField(); i++ {
//		field := cfgType.Field(i)
//		name := strcase.ToSnake(field.Name)
//		argName := fmt.Sprintf("%s_%s", moduleGlobalName, name)
//		switch field.Type.Kind() {
//		case reflect.Int:
//			a.addFlagInt64(a.cmdBackup, argName, cfgValue.Field(i).Int())
//		case reflect.String:
//			a.addFlagString(a.cmdBackup, argName, cfgValue.Field(i).String())
//		case reflect.Bool:
//			a.addFlagBool(a.cmdBackup, argName, cfgValue.Field(i).Bool())
//		case reflect.Struct:
//			panic(fmt.Errorf("module %s config contains struct field %s", moduleGlobalName, field.Name))
//		}
//	}
//	return nil
//}
