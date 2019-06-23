package common

type ConfigBackup struct {
	Connect   *ConfigModule   `yaml:'connect'`
	Input     *ConfigModule   `yaml:'input'`
	Compress  *ConfigModule   `yaml:'compress'`
	Encrypt   *ConfigModule   `yaml:'encrypt'`
	Outputs   []*ConfigModule `yaml:'outputs'`
	Notifiers []*ConfigModule `yaml:'notifiers'`
}
