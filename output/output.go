package output

import (
	"github.com/copybird/copybird/core"
)

type Output interface {
	core.PipeComponent
	InitOutput() error
}
