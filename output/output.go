package output

import (
	"github.com/copybird/copybird/core"
)

type Output interface {
	core.PipeComponent
	Upload() error
}
