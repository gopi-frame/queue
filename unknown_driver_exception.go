package queue

import (
	"fmt"

	"github.com/gopi-frame/exception"
)

type UnknownDriverException struct {
	*exception.Exception
	driver string
}

func (e *UnknownDriverException) Driver() string {
	return e.driver
}

func NewUnknownDriverException(driver string) *UnknownDriverException {
	exp := new(UnknownDriverException)
	exp.driver = driver
	exp.Exception = exception.NewException(fmt.Sprintf("queue: unknown driver \"%s\"", driver))
	return exp
}
