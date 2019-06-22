package core

import (
	"fmt"
)

type ModuleError struct {
	Module Module
	Err    error
}

func (me ModuleError) Error() string {
	return fmt.Sprintf("module %s err: %s", me.Module.GetName(), me.Err)
}
