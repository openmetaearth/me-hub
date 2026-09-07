package utest

import (
	"errors"
	"strings"

	errorsmod "cosmossdk.io/errors"
)

type Truer interface {
	True(value bool, msgAndArgs ...interface{})
}

func IsErr(t Truer, actual, expected error) {
	ok := actual != nil && expected != nil && (
		errorsmod.IsOf(actual, expected) || errorsmod.IsOf(expected, actual) ||
			errors.Is(actual, expected) || errors.Is(expected, actual) ||
			strings.Contains(actual.Error(), expected.Error()) || strings.Contains(expected.Error(), actual.Error()))
	t.True(ok, `error is not an instance of expected: expected: %T, %s: got: %T, %s`, expected, expected, actual, actual)
}
