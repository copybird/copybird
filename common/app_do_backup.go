package common

import (
	"errors"
	"fmt"
	"github.com/copybird/copybird/core"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

func (a *App) DoBackup() error {
	lFlags := a.cmdBackup.LocalFlags()
	mInputArgs := lFlags.Lookup("input")
	if mInputArgs == nil {
		mInputArgs = lFlags.ShorthandLookup("i")
		if mInputArgs == nil {
			return errors.New("need input")
		}
	}
	mOutputArgs := lFlags.Lookup("output")
	if mOutputArgs == nil {
		mOutputArgs = lFlags.ShorthandLookup("o")
		if mOutputArgs == nil {
			return errors.New("need at least one output")
		}
	}
	mInputName, mInputParams := parseArgs(mInputArgs.Value.String())
	mOutputName, mOutputParams := parseArgs(mOutputArgs.Value.String())
	log.Printf("input: %s %+v", mInputName, mInputParams)
	log.Printf("output: %s %+v", mOutputName, mOutputParams)

	mInput := core.GetModule(core.ModuleGroupBackup, core.ModuleTypeInput, mInputName)
	if mInput == nil {
		return fmt.Errorf("input module %s not found", mInputName)
	}

	mOutput := core.GetModule(core.ModuleGroupBackup, core.ModuleTypeOutput, mOutputName)
	if mOutput == nil {
		return fmt.Errorf("output module %s not found", mInputName)
	}

	mInput.InitModule(mInput.GetConfig())
	mOutput.InitModule(mOutput.GetConfig())

	wg := sync.WaitGroup{}
	wg.Add(2)

	chanError := make(chan error, 1000)

	inputReader, inputWriter := io.Pipe()
	mInput.InitPipe(inputWriter, nil)
	mOutput.InitPipe(nil, inputReader)

	go runModule(mInput, &wg, chanError)
	go runModule(mOutput, &wg, chanError)

	wg.Wait()

	for {
		err, ok := <-chanError
		if !ok {
			break
		}
		log.Printf("err: %s", err)
	}

	return nil
}

func runModule(module core.Module, wg *sync.WaitGroup, chanError chan error) {
	defer func(t time.Time) {
		wg.Done()
		if err := recover(); err != nil {
			chanError <- fmt.Errorf("module %s/%s panic: %s", module.GetType(), module.GetName(), err)
		}
		log.Printf("module %s/%s done by %.2fms", module.GetType(), module.GetName(), time.Since(t).Seconds()*1000)
	}(time.Now())
	err := module.Run()
	if err != nil {
		chanError <- err
	}
}

func parseArgs(args string) (string, map[string]string) {
	var moduleName string
	var moduleParams = make(map[string]string)

	parts := strings.Split(args, "::")
	moduleName = parts[0]
	if len(parts) > 1 {
		for _, param := range parts[1:] {
			paramParts := strings.Split(param, "=")
			moduleParams[paramParts[0]] = paramParts[1]
		}
	}

	return moduleName, moduleParams
}
