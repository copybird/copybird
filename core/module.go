package core

import (
	"context"
	"io"
)

type Module interface {
	GetName() string
	GetGroup() ModuleGroup
	GetType() ModuleType
	GetConfig() interface{}
	InitModule(cfg interface{}) error
	Run(ctx context.Context) error
	Close() error
	InitPipe(w io.Writer, r io.Reader) error
}
