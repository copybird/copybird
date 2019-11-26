package common

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/copybird/copybird/core"
	"github.com/iancoleman/strcase"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"
)

const (
	envInput    = "COPYBIRD_INPUT"
	envOutput   = "COPYBIRD_OUTPUT"
	envCompress = "COPYBIRD_COMPRESS"
	envEncrypt  = "COPYBIRD_ENCRYPT"
)

func (a *App) DoBackup() error {
	var mInput, mCompress, mEncrypt, mOutput core.Module
	var err error
	// Read CLI args
	lFlags := a.cmdBackup.LocalFlags()
	mInputArgs := lFlags.Lookup("input")
	if mInputArgs == nil {
		mInputArgs = lFlags.ShorthandLookup("i")
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
	}

	// Load modules
	if mInputArgs != nil && mInputArgs.Value.String() != "" {
		mInput, err = loadModuleFromArgs(core.ModuleGroupBackup, core.ModuleTypeInput, mInputArgs.Value.String())
		if err != nil {
			return err
		}
	} else {
		mInput, err = loadModuleFromEnv(core.ModuleGroupBackup, core.ModuleTypeInput, envInput)
		if err != nil {
			return fmt.Errorf("can't get input module config: %s", err)
		}
	}

	if mOutputArgs != nil && mOutputArgs.Value.String() != "" {
		mOutput, err = loadModuleFromArgs(core.ModuleGroupBackup, core.ModuleTypeOutput, mOutputArgs.Value.String())
		if err != nil {
			return err
		}
	} else {
		mInput, err = loadModuleFromEnv(core.ModuleGroupBackup, core.ModuleTypeOutput, envOutput)
		if err != nil {
			return fmt.Errorf("can't get output module config: %s", err)
		}
	}

	if mCompressArgs != nil && mCompressArgs.Value.String() != "" {
		mCompress, err = loadModuleFromArgs(core.ModuleGroupBackup, core.ModuleTypeCompress, mCompressArgs.Value.String())
		if err != nil {
			return err
		}
	} else {
		mCompress, err = loadModuleFromEnv(core.ModuleGroupBackup, core.ModuleTypeOutput, envCompress)
	}

	if mEncryptArgs != nil && mEncryptArgs.Value.String() != "" {
		mEncrypt, err = loadModuleFromArgs(core.ModuleGroupBackup, core.ModuleTypeEncrypt, mEncryptArgs.Value.String())
		if err != nil {
			return err
		}
	} else {
		mEncrypt, err = loadModuleFromEnv(core.ModuleGroupBackup, core.ModuleTypeOutput, envEncrypt)
	}

	// Run pipeline
	eg, ctx := errgroup.WithContext(context.Background())
	nextReader, nextWriter := io.Pipe()

	eg.Go(runModule(ctx, mInput, nextWriter, nil))
	if mCompress != nil {
		interimReader, interimWriter := io.Pipe()
		eg.Go(runModule(ctx, mCompress, interimWriter, nextReader))
		nextReader = interimReader
	}
	if mEncrypt != nil {
		interimReader, interimWriter := io.Pipe()
		eg.Go(runModule(ctx, mEncrypt, interimWriter, nextReader))
		nextReader = interimReader
	}
	eg.Go(runModule(ctx, mOutput, nil, nextReader))

	return eg.Wait()
}

func loadModuleFromArgs(mGroup core.ModuleGroup, mType core.ModuleType, args string) (core.Module, error) {
	name, params := parseArgs(args)
	module := core.GetModule(mGroup, mType, name)
	if module == nil {
		return nil, fmt.Errorf("module %s/%s not found", mType, name)
	}
	config := module.GetConfig()
	loadConfig(config, params)
	log.Printf("module %s/%s config: %+v", mType, name, config)
	if err := module.InitModule(config); err != nil {
		return nil, fmt.Errorf("init module %s/%s err: %s", mType, name, err)
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

func runModule(ctx context.Context, module core.Module, writer io.WriteCloser, reader io.ReadCloser) func() error {
	return func() error {
		defer func() {
			if writer != nil {
				writer.Close()
			}
			if reader != nil {
				reader.Close()
			}
		}()
		t := time.Now()
		err := module.InitPipe(writer, reader)
		if err != nil {
			return fmt.Errorf("module %s/%s pipe initialization err: %s", module.GetType(), module.GetName(), err)
		}
		err = module.Run(ctx)
		if err != nil {
			return fmt.Errorf("module %s/%s execution err: %s", module.GetType(), module.GetName(), err)
		}
		log.Printf("module %s/%s done by %.2fms", module.GetType(), module.GetName(), time.Since(t).Seconds()*1000)
		return nil
	}
}

func parseArgs(args string) (string, map[string]string) {
	var moduleName string
	var moduleParams = make(map[string]string)

	parts := strings.Split(args, "::")
	if len(parts) == 0 {
		return "", nil
	}
	moduleName = parts[0]
	for _, param := range parts[1:] {
		paramParts := strings.Split(param, "=")
		if len(paramParts) == 2 && paramParts[0] != "" {
			moduleParams[paramParts[0]] = paramParts[1]
		}
	}

	return moduleName, moduleParams
}

func loadModuleFromEnv(mGroup core.ModuleGroup, mType core.ModuleType, envPrefix string) (core.Module, error) {
	mName, defined := os.LookupEnv(envPrefix)
	if !defined {
		return nil, fmt.Errorf("%s/%s component not defined", mGroup, mType)
	}
	module := core.GetModule(mGroup, mType, mName)
	if module == nil {
		return nil, fmt.Errorf("module %s/%s not found", mType, mName)
	}
	config := module.GetConfig()
	err := envconfig.Process(envPrefix, config)
	if err != nil {
		return nil, fmt.Errorf("module %s/%s env parse err: %s", mType, mName, err)
	}
	log.Printf("module %s/%s config: %+v", mType, mName, config)
	if err := module.InitModule(config); err != nil {
		return nil, fmt.Errorf("init module %s/%s err: %s", mType, mName, err)
	}
	return module, nil
}
