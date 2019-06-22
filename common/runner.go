package common

import (
	"github.com/copybird/copybird/core"
)

type Runner struct {
	moduleInput core.Module
	moduleCompress core.Module
	moduleEncrypt core.Module
	moduleOutput core.Module
	moduleNotifiers []core.Module
}

func (r *Runner) Run() {
	wg := sync.WaitGroup{}
	// for input and output
	wg.Add(2) 
	// for compress
	if moduleCompress != nil {
		wg.Add(1) 
	}
	// for encryption
	if moduleEncrypt != nil {
		wg.Add(1) 
	}
	chanError := make(chan error, 1000)
	nextReader, nextWriter := io.Pipe()
	go runModule(r.moduleInput, nextWriter, nil, &wg, chanError)
	if r.moduleCompress != nil {
		_nextReader, _nextWriter := io.Pipe()
		go runModule(r.moduleCompress, _nextWriter, nextReader, &wg, chanError)
		nextReader = _nextReader
	}
	if r.moduleEncrypt != nil {
		_nextReader, _nextWriter := io.Pipe()
		go runModule(r.moduleEncrypt, _nextWriter, nextReader, &wg, chanError)
		nextReader = _nextReader
	}
	_nextReader, _nextWriter := io.Pipe()
	go runModule(r.moduleOutput, nil, nextReader, &wg, chanError)
	wg.Wait()
	for {
		err, ok := <- chanError
		if !ok {
			break
		}
		log.Printf("")
	}
}

func (r *Runner) runModule(module core.Module, writer io.Writer, reader io.Reader, wg *sync.WaitGroup, chanError chan error) {
	defer wg.Done()
	err := module.InitPipe(writer, reader)
	if err != nil {
		chanError <- &ModuleError{Module: module, Err: err}
		return
	}
	err = module.Run()
	if err != nil {
		chanError <- &ModuleError{Module: module, Err: err}
	}
}