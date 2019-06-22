package output

import (
	"github.com/copybird/copybird/core"
)

type Output interface {
	core.Module
	core.PipeComponent
}
 