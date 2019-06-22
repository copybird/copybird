package common

import (
	"fmt"
	compress_gzip "github.com/copybird/copybird/compress/gzip"
	connect_ssh "github.com/copybird/copybird/connect/ssh"
	"github.com/copybird/copybird/core"
	"github.com/copybird/copybird/encryption/aesgcm"
	"github.com/copybird/copybird/input/mysql"
	"github.com/copybird/copybird/notifier/awsses"
	"github.com/copybird/copybird/notifier/email"
	"github.com/copybird/copybird/notifier/kafka"
	"github.com/copybird/copybird/notifier/nats"
	"github.com/copybird/copybird/notifier/pagerduty"
	"github.com/copybird/copybird/notifier/pushbullet"
	"github.com/copybird/copybird/notifier/rabbitmq"
	"github.com/copybird/copybird/notifier/slack"
	"github.com/copybird/copybird/notifier/telegram"
	"github.com/copybird/copybird/notifier/twillio"
	"github.com/copybird/copybird/notifier/webcallback"
	"github.com/copybird/copybird/output/gcp"
	"github.com/copybird/copybird/output/http"
	"github.com/copybird/copybird/output/local"
	"github.com/copybird/copybird/output/s3"
	"github.com/copybird/copybird/output/scp"
	"github.com/spf13/cobra"
	"strings"

	//"strings"

	//"log"
	"reflect"
	//"strings"

	"github.com/iancoleman/strcase"

	// lz4_compress "github.com/copybird/copybird/compress/lz4"
	// lz4_decompress "github.com/copybird/copybird/decompress/lz4"
)

type ModuleType int

const (
	ModuleTypeConnect ModuleType = iota
	ModuleTypeInput
	ModuleTypeCompress
	ModuleTypeEncryption
	ModuleTypeOutput
	ModuleTypeNotify
)

func (m ModuleType) String() string {
	return [...]string{
		"connect",
		"input",
		"compress",
		"encryption",
		"output",
		"notify",
	}[m]
}

func (a *App) Setup() error {
	a.addFlagString(a.cmdBackup, "input", "mysql")
	a.addFlagString(a.cmdBackup, "compress", "")
	a.addFlagString(a.cmdBackup, "encrypt", "")
	a.addFlagString(a.cmdBackup, "output", "local")
	a.addFlagString(a.cmdBackup, "notify", "slack")
	a.RegisterModules()
	return nil
}

func (a *App) addFlagString(cmd *cobra.Command, name string, defaultValue string) {
	a.vars[name] = cmd.Flags().String(name, defaultValue, fmt.Sprintf("env %s", strings.ToUpper(name)))
}

func (a *App) addFlagInt64(cmd *cobra.Command, name string, defaultValue int64) {
	a.vars[name] = cmd.Flags().Int64(name, defaultValue, fmt.Sprintf("env %s", strings.ToUpper(name)))
}

func (a *App) addFlagBool(cmd *cobra.Command, name string, defaultValue bool) {
	a.vars[name] = cmd.Flags().Bool(name, defaultValue, fmt.Sprintf("env %s", strings.ToUpper(name)))
}

func (a *App) RegisterModules() {
	a.RegisterModule(ModuleTypeConnect, &connect_ssh.Ssh{})

	a.RegisterModule(ModuleTypeInput, &mysql.MySQLDumper{})

	a.RegisterModule(ModuleTypeCompress, &compress_gzip.Compress{})

	a.RegisterModule(ModuleTypeEncryption, &aesgcm.EncryptionAESGCM{})

	a.RegisterModule(ModuleTypeOutput, &local.Local{})
	a.RegisterModule(ModuleTypeOutput, &s3.S3{})
	a.RegisterModule(ModuleTypeOutput, &gcp.GCP{})
	a.RegisterModule(ModuleTypeOutput, &http.Http{})
	a.RegisterModule(ModuleTypeOutput, &scp.SCP{})

	a.RegisterModule(ModuleTypeNotify, &awsses.AwsSes{})
	a.RegisterModule(ModuleTypeNotify, &email.Email{})
	a.RegisterModule(ModuleTypeNotify, &kafka.Kafka{})
	a.RegisterModule(ModuleTypeNotify, &nats.Nats{})
	a.RegisterModule(ModuleTypeNotify, &pagerduty.PagerDuty{})
	a.RegisterModule(ModuleTypeNotify, &pushbullet.Local{})
	a.RegisterModule(ModuleTypeNotify, &rabbitmq.RabbitMQ{})
	a.RegisterModule(ModuleTypeNotify, &slack.Local{})
	a.RegisterModule(ModuleTypeNotify, &telegram.Local{})
	a.RegisterModule(ModuleTypeNotify, &twillio.Twillio{})
	a.RegisterModule(ModuleTypeNotify, &webcallback.Callback{})
}

func (a *App) RegisterModule(moduleType ModuleType, module core.Module) error {
	moduleGlobalName := fmt.Sprintf("%s_%s", moduleType.String(), module.GetName())
	a.registeredModules[moduleGlobalName] = module
	cfg := module.GetConfig()
	cfgValue := reflect.Indirect(reflect.ValueOf(cfg))
	cfgType := cfgValue.Type()
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		name := strcase.ToSnake(field.Name)
		argName := fmt.Sprintf("%s_%s", moduleGlobalName, name)
		switch (field.Type.Kind()) {
		case reflect.Int:
			a.addFlagInt64(a.cmdBackup, argName, cfgValue.Field(i).Int())
		case reflect.String:
			a.addFlagString(a.cmdBackup, argName, cfgValue.Field(i).String())
		case reflect.Bool:
			a.addFlagBool(a.cmdBackup, argName, cfgValue.Field(i).Bool())
		}
	}
	return nil
}
