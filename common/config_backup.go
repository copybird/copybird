package common

type ConfigBackup struct {
	Connect   ConfigModule
	Input     ConfigModule
	Compress  ConfigModule
	Encrypt   ConfigModule
	Outputs   []ConfigModule
	Notifiers []ConfigModule
}
