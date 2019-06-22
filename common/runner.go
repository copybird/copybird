package common

import (
	"github.com/copybird/copybird/core"
)

type Runner struct {
	moduleInput core.Module
	moduleCompress core.Module
	moduleEncrypt core.Module
	moduleOutputs []core.Module
}