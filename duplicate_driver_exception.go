package queue

import (
	"fmt"

	"github.com/gopi-frame/exception"
)

type DuplicateDriverException struct {
	*exception.Exception
	driver string
}

func (e *DuplicateDriverException) Driver() string {
	return e.driver
}

func NewDuplicateDriverException(driver string) *DuplicateDriverException {
	exp := new(DuplicateDriverException)
	exp.driver = driver
	exp.Exception = exception.NewException(fmt.Sprintf("sql: Register called twice for driver \"%s\"", driver))
	return exp
}
