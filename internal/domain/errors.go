package domain

import (
	"fmt"
)

type Code string

const (
	ErrorCodeTeamExists  Code = "TEAM_EXISTS"
	ErrorCodePRExists    Code = "PR_EXISTS"
	ErrorCodePRMerged    Code = "PR_MERGED"
	ErrorCodeNotAssigned Code = "NOT_ASSIGNED"
	ErrorCodeNoCandidate Code = "NO_CANDIDATE"
	ErrorCodeNotFound    Code = "NOT_FOUND"
)

type Error struct {
	Code    Code
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code Code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}
