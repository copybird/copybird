package common

type ConfigRestore struct {
	Connect    ConfigModule
	Input      ConfigModule
	Decompress ConfigModule
	Decrypt    ConfigModule
	Output     ConfigModule
	Notifiers  []ConfigModule
}
