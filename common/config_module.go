package common

type ConfigModule struct {
	Type   string      `yaml:'type'`
	Config interface{} `yaml:'config'`
}
