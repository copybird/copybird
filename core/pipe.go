package core

import (
	"io"
)

type PipeComponent interface {
	Init(w io.Writer, r io.Reader) error
	Run() error
	Close() error
}
