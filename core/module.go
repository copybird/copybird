package core

type Module interface {
	GetName() string
	GetConfig() interface{}
	InitModule(cfg interface{}) error
	Run() error
	Close() error
}
