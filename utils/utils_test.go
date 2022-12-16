package utils

import (
	"errors"
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

func TestStack(t *testing.T) {
	s := "Major problem:"

	errs := []error{
		errors.New("1st err"),
		errors.New("2nd err"),
		errors.New("3rd err"),
		errors.New("4th err"),
	}

	err := Stack(errs, errors.New(s))
	assert.Equal(t, s+" 1st err 2nd err 3rd err 4th err", err.Error(), "Error strings should be equal")

	errs = []error{}
	err = Stack(errs, errors.New(s))
	assert.Equal(t, s, err.Error(), "Error strings should be equal")
}
