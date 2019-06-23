package common

type BackupConfig struct {
	Connect   ModuleConfig
	Input     ModuleConfig
	Compress  ModuleConfig
	Encrypt   ModuleConfig
	Outputs   []ModuleConfig
	Notifiers []ModuleConfig
}

