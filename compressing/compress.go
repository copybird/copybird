package compressing

import "github.com/copybird/copybird/core"

type Output interface {
	core.PipeComponent
	InitCompress(i int) error
}