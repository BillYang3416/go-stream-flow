package apperrors

import (
	"errors"
	"fmt"
)

type UniqueConstraintError struct {
	Message        string
	LoggingContext string
}

func (e *UniqueConstraintError) Error() string {
	return e.Message
}

func NewUniqueConstraintError(msg string, loggingContext string, args ...interface{}) *UniqueConstraintError {
	return &UniqueConstraintError{Message: fmt.Sprintf(msg, args...), LoggingContext: loggingContext}
}

func IsUniqueConstraintError(err error) bool {
	var uce *UniqueConstraintError
	return errors.As(err, &uce)
}

func AsUniqueConstraintError(err error) (*UniqueConstraintError, bool) {
	var uce *UniqueConstraintError
	ok := errors.As(err, &uce)
	return uce, ok
}

type NoRowsAffectedError struct {
	Message        string
	LoggingContext string
}

func (e *NoRowsAffectedError) Error() string {
	return e.Message
}

func NewNoRowsAffectedError(msg string, loggingContext string, args ...interface{}) *NoRowsAffectedError {
	return &NoRowsAffectedError{Message: fmt.Sprintf(msg, args...), LoggingContext: loggingContext}
}

func IsNoRowsAffectedError(err error) bool {
	var nrae *NoRowsAffectedError
	return errors.As(err, &nrae)
}

func AsNoRowsAffectedError(err error) (*NoRowsAffectedError, bool) {
	var nrae *NoRowsAffectedError
	ok := errors.As(err, &nrae)
	return nrae, ok
}
