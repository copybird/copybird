package encryption

import (
	"github.com/copybird/copybird/core"
)

type Encryption interface {
	core.PipeComponent
	InitEncryption(key []byte) error
}
