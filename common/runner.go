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

func (r *Runner) Run() {
	wg := sync.WaitGroup{}
	// for input and outputs
	wg.Add(1 + len(r.moduleOutputs)) 
	// for compress
	if moduleCompress != nil {
		wg.Add(1) 
	}
	// for encryption
	if moduleEncrypt != nil {
		wg.Add(1) 
	}
}