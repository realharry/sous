package sous

import (
	"fmt"
	"testing"

	"github.com/nyarly/testify/assert"
	"github.com/pkg/errors"
)

func TestIsTransientResolveError(t *testing.T) {
	assert := assert.New(t)

	assert.False(IsTransientResolveError(fmt.Errorf("Hi!")))
	assert.True(IsTransientResolveError(&CreateError{}))
	assert.True(IsTransientResolveError(errors.Wrap(&CreateError{}, "even if wrapped")))
}
