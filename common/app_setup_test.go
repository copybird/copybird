package common

import (
	"testing"

	"gotest.tools/assert"
)

func TestAppSetup(t *testing.T) {
	app := NewApp()
	assert.NilError(t, app.Setup())
}