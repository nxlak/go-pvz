package errs

import (
	"fmt"
	"strings"
)

type Field struct {
	Key   string
	Value any
}

type AppError struct {
	Code    string
	Message string
	Err     error
	Fields  []Field
}

func (e *AppError) Error() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("code = %s, message = %s", e.Code, e.Message))

	for _, f := range e.Fields {
		b.WriteString(fmt.Sprintf(" %s = %v", f.Key, f.Value))
	}

	if e.Err != nil {
		b.WriteString(fmt.Sprintf(" | cause = %s", e.Err))
	}

	return b.String()
}

func New(code, message string, fields ...any) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Fields:  parseFields(fields),
	}
}

func Wrap(code, message string, err error, fields ...any) *AppError {
	return &AppError{
		Code:    code,
		Err:     err,
		Message: message,
		Fields:  parseFields(fields),
	}
}
