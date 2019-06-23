package input

import "github.com/copybird/copybird/core"

// BackupDatasource interface is basic interface for all inputs (databases, files)
type BackupDatasource interface {
	core.Module
	core.PipeComponent
}
