package common

type RestoreConfig struct {
	Connect    ModuleConfig
	Input      ModuleConfig
	Decompress ModuleConfig
	Decrypt    ModuleConfig
	Output     ModuleConfig
	Notifiers  []ModuleConfig
}
