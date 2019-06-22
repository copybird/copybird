package common

import (
	"testing"

	"gotest.tools/assert"
)

func TestAppSetup(t *testing.T) {
	app := App{}
	assert.NilError(t, app.Setup())
}