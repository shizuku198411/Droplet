package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileSpecLoader(t *testing.T) {
	// == arrange ==

	// == act ==
	fileSpecLoader := newFileSpecLoader()

	// == assert ==
	assert.NotNil(t, fileSpecLoader)
}
