package genericpam

import (
	"errors"
	"fmt"
)

func (e PamError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func toPamError(err error) *PamError {
	pamError := &PamError{}
	if !errors.As(err, pamError) {
		pamError.Code = PAMERRUNDEFINED
		pamError.Message = err.Error()
	}
	return pamError
}
