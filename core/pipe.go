package core

import (
	"io"
)

type PipeComponent interface {
	InitPipe(w io.Writer, r io.Reader) error
	Run() error
	Close() error
}
