// Package utils provides utility functions used throughout pam.
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"net"
)

var rndSize = big.NewInt(1 << 32)

func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func RandomInt() int {
	n, _ := rand.Int(rand.Reader, rndSize)
	return int(n.Int64())
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
func GetFreePort() (int, string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, "", err
	}

	defer func() {
		_ = l.Close()
	}()
	return l.Addr().(*net.TCPAddr).Port, l.Addr().String(), nil
}
