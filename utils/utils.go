// Package utils provides utility functions used throughout pam.
package utils

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomInt() int {
	return rand.Int()
}

// Stack combines an array of errors into a single error via formatting.
func Stack(errs []error, target error) error {
	if len(errs) == 0 {
		return target
	}

	return Stack(errs[1:], fmt.Errorf("%s %w", target, errs[0]))
}

// Ptr returns the pointer to an argument, useful for string literals.
func Ptr[T any](t T) *T {
	return &t
}

// OrZeroValue returns the value referenced by the pointer argument ptr, or if nil it returns the zero value.
// For example, OrZeroValue[string](nil) returns the zero value for string ("").
func OrZeroValue[T any](ptr *T) T {
	var res T
	if ptr != nil {
		res = *ptr
	}
	return res
}

// OrDefault returns the value referenced by the pointer argument ptr, or if nil it returns the default value.
func OrDefault[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}

// GetFreePort returns a free open port that is ready to use.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = l.Close()
	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}
