package output

import (
	"github.com/copybird/copybird/core"
)

type Output interface {
	core.PipeComponent
	InitOutput(map[string]string) error
}
