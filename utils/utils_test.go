package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrZeroValue(t *testing.T) {
	assert.Zero(t, OrZeroValue[string](nil))
	assert.Zero(t, OrZeroValue[int](nil))
	assert.Zero(t, OrZeroValue[bool](nil))
}

func TestPtr(t *testing.T) {
	var s = "str"

	assert.Equal(t, &s, Ptr(s))
	assert.NotNil(t, Ptr("literal string"))
}
