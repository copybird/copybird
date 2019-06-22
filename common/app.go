package common

import (
	"github.com/copybird/copybird/core"
	"github.com/spf13/cobra"
	"io"
	"log"
	"sync"

	//"log"

	//"github.com/spf13/cobra"
)

type App struct {
	registeredModules map[string]core.Module
	cmmRoot           *cobra.Command
	cmdBackup         *cobra.Command
	vars              map[string]interface{}
}

func NewApp() *App {
	return &App{
		registeredModules: make(map[string]core.Module),
		vars:              make(map[string]interface{}),
	}
}

func (a *App) Run() error {

	var rootCmd = &cobra.Command{Use: "copybird"}
	a.cmdBackup = &cobra.Command{
		Use:   "backup",
		Short: "Start new backup",
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdCallback(a.DoBackup),
	}
	rootCmd.AddCommand(a.cmdBackup)
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

func (a *App) DoBackup() error {
	moduleNameInput := *a.vars["input"].(*string)
	moduleNameCompress := *a.vars["compress"].(*string)
	moduleNameEncrypt := *a.vars["encrypt"].(*string)
	moduleNameOutput := *a.vars["output"].(*string)
	log.Printf("module input: %s", moduleNameInput)
	if moduleNameCompress != "" {
		log.Printf("module compress: %s", moduleNameInput)
	}
	if moduleNameEncrypt != "" {
		log.Printf("module compress: %s", moduleNameEncrypt)
	}
	var err error
	log.Printf("module output: %s", moduleNameOutput)
	moduleInput := a.getModule("input_" + moduleNameInput)
	moduleInputConfig := a.getModuleConfig(ModuleTypeInput, moduleInput)
	if err := moduleInput.InitModule(moduleInputConfig); err != nil {
		return err
	}
	moduleOutput := a.getModule("output_" + moduleNameOutput)
	moduleOutputConfig := a.getModuleConfig(ModuleTypeOutput, moduleOutput)
	if err := moduleOutput.InitModule(moduleOutputConfig); err != nil {
		return err
	}
	inputReader, inputWriter := io.Pipe()
	moduleInput.InitPipe(inputWriter, nil)
	moduleOutput.InitPipe(nil, inputReader)
	wg := sync.WaitGroup{}
	wg.Add(2)
	chanError := make(chan error, 1000)
	go func(wg *sync.WaitGroup, chErr chan error) {
		defer wg.Done()
		if err := moduleInput.Run(); err != nil {
			chanError <- err
		}
	}(&wg, chanError)
	wg.Wait()
	err, ok := <-chanError
	if !ok {
		return err
	}
	return nil
}

func (a *App) getModule(name string) core.Module {
	return a.registeredModules[name]
}

func (a *App) getModuleConfig(moduleType ModuleType, module core.Module) interface{} {
	//moduleGlobalName := fmt.Sprintf("%s_%s", moduleType.String(), module.GetName())
	cfg := module.GetConfig()
	//cfgValue := reflect.Indirect(reflect.ValueOf(cfg))
	//cfgType := cfgValue.Type()
	//for i := 0; i < cfgValue.NumField(); i++ {
	//	field := cfgType.Field(i)
	//	name := strcase.ToSnake(field.Name)
	//	argName := fmt.Sprintf("%s_%s", moduleGlobalName, name)
	//	argValue := reflect.Indirect(reflect.ValueOf(a.vars[argName]))
	//	cfgValue.Field(i).Set(reflect.ValueOf(argValue.String()))
	//}
	log.Printf("config %s: %#v", module.GetName(), cfg)
	return cfg
}
