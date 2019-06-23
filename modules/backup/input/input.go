package input

import "github.com/copybird/copybird/core"

// Input interface is basic interface for all inputs (databases, files)
type Input interface {
	core.Module
}
