package common

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/copybird/copybird/core"
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

	mCompressArgs := lFlags.Lookup("compress")
	if mCompressArgs == nil {
		mCompressArgs = lFlags.ShorthandLookup("z")
	}

	mEncryptArgs := lFlags.Lookup("encrypt")
	if mEncryptArgs == nil {
		mEncryptArgs = lFlags.ShorthandLookup("e")
	}

	mOutputArgs := lFlags.Lookup("output")
	if mOutputArgs == nil {
		mOutputArgs = lFlags.ShorthandLookup("o")
		if mOutputArgs == nil {
			return errors.New("need at least one output")
		}
	}

	mInput, err := loadModule(core.ModuleGroupBackup, core.ModuleTypeInput, mInputArgs.Value.String())
	if err != nil {
		return err
	}
	mOutput, err := loadModule(core.ModuleGroupBackup, core.ModuleTypeOutput, mOutputArgs.Value.String())
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	nextReader, nextWriter := io.Pipe()

	go runModule(mInput, nextWriter, nil, &wg)

	if mCompressArgs != nil && mCompressArgs.Value.String() != "" {
		mCompress, err := loadModule(core.ModuleGroupBackup, core.ModuleTypeCompress, mCompressArgs.Value.String())
		if err != nil {
			return err
		}
		_nextReader, _nextWriter := io.Pipe()
		wg.Add(1)
		go runModule(mCompress, _nextWriter, nextReader, &wg)
		nextReader = _nextReader
	}

	if mEncryptArgs != nil && mEncryptArgs.Value.String() != "" {
		mEncrypt, err := loadModule(core.ModuleGroupBackup, core.ModuleTypeEncrypt, mEncryptArgs.Value.String())
		if err != nil {
			return err
		}
		_nextReader, _nextWriter := io.Pipe()
		wg.Add(1)
		go runModule(mEncrypt, _nextWriter, nextReader, &wg)
		nextReader = _nextReader
	}

	go runModule(mOutput, nil, nextReader, &wg)

	wg.Wait()

	return nil
}

func loadModule(mGroup core.ModuleGroup, mType core.ModuleType, args string) (core.Module, error) {
	name, params := parseArgs(args)
	module := core.GetModule(mGroup, mType, name)
	if module == nil {
		return nil, fmt.Errorf("module %s not found", name)
	}
	config := module.GetConfig()
	loadConfig(config, params)
	log.Printf("module %s/%s config: %+v", mType, name, config)
	if err := module.InitModule(config); err != nil {
		return nil, err
	}
	return module, nil
}

func loadConfig(cfg interface{}, params map[string]string) error {
	cfgValue := reflect.Indirect(reflect.ValueOf(cfg))
	cfgType := cfgValue.Type()

	for pName, pValue := range params {
		for i := 0; i < cfgType.NumField(); i++ {
			fieldValue := cfgValue.Field(i)
			fieldType := cfgType.Field(i)
			if strcase.ToSnake(fieldType.Name) == pName {
				switch fieldType.Type.Kind() {
				case reflect.String:
					fieldValue.SetString(pValue)
				case reflect.Int:
					val, err := strconv.ParseInt(pValue, 10, 63)
					if err != nil {
						return err
					}
					fieldValue.SetInt(val)
				case reflect.Bool:
					val, err := strconv.ParseBool(pValue)
					if err != nil {
						return err
					}
					fieldValue.SetBool(val)
				default:
					return fmt.Errorf("unsupported config param type: %s %s",
						pName,
						fieldType.Type.Kind().String())
				}
			}
		}
	}
	return nil
}

func runModule(module core.Module, writer io.WriteCloser, reader io.ReadCloser, wg *sync.WaitGroup) {
	defer func(t time.Time) {
		if writer != nil {
			writer.Close()
		}
		if reader != nil {
			reader.Close()
		}
		wg.Done()
		if err := recover(); err != nil {
			log.Printf("module %s/%s err: %s", module.GetType(), module.GetName(), err)
		}
		log.Printf("module %s/%s done by %.2fms", module.GetType(), module.GetName(), time.Since(t).Seconds()*1000)
	}(time.Now())
	err := module.InitPipe(writer, reader)
	if err != nil {
		panic(err)
	}
	err = module.Run()
	if err != nil {
		panic(err)
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
