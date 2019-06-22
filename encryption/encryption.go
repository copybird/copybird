package encryption

import (
	"github.com/copybird/copybird/core"
)

type Encryption interface {
	core.Module
	core.PipeComponent
}
