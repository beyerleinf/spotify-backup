package assert

import (
	"fmt"
	"runtime"
)

// Assert checks if the given condition is true.
// If it's false, it panics with a message including the file and line number.
func Assert(condition bool, message string) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("Assertion failed at %s:%d: %s", file, line, message))
	}
}

// Equal checks if two values are equal.
// If they're not, it panics with a message including the file and line number.
func Equal(expected, actual interface{}, message string) {
	if expected != actual {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("Assertion failed at %s:%d: %s. Expected \"%v\", got \"%v\"", file, line, message, expected, actual))
	}
}

// NotEqual checks if two values are not equal.
// If they are, it panics with a message including the file and line number.
func NotEqual(expected, actual interface{}, message string) {
	if expected == actual {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("Assertion failed at %s:%d: %s. Expected \"%v\" not to equal \"%v\"", file, line, message, actual, expected))
	}
}

// NotNil checks if a value is not nil.
// If it is nil, it panics with a message including the file and line number.
func NotNil(value interface{}, message string) {
	if value == nil {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("Assertion failed at %s:%d: %s. Value is nil", file, line, message))
	}
}
