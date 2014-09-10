/* package error provides utilities for creating and reporting errors.

Currently supported ErrorCodes are

	Configuration
	Library
	MissingFile
	Sanity
	Value

which correspond to
configuration errors,
library errors,
missing file errors,
sanity errors,
and value errors,
respectively.

Configuration errors indicate that the user has given the program invalid
exernal input and, most importantly, can fix the error entirely through
external input. Configuration errors must not be returned after the
initialization period of the program. All configuration varibles with invalid
value ranges must result in configuration errors when encountering these
value.

Library errors indicate that some exernal library (most likely a curses
library) has thrown an error that the program can't recover from. All such
library errors must be reported.

Missing file errors indicate that a required file does not exist or was
placed in the wrong location. Any function which reads from a file must
result in a file error unless that function can recover internally. For
modularity purposes it's probably best that no function can recover
internally.

Sanity errors indicate that something impossible has occured. This could mean
that there was a non-exhaustive switch statement or that some data structure
invariant was not upheld, or something similar. Sanity errors must be used in
all places where adding functionality requires updating multiple parts of the
code, but are otherwise up to the programmer's discretion.

Value errors indicate that a function has been given a parameter which is
outside its valid value range. All of a package's externally visible
functions must give value errors when encountering such a value. Internal
package functions may return these errors at the programmer's discretion. If
input is in the form of labeled integers (like ErrorCodes) and there is no
other potentially erroneous input, value errors don't need to be returned.*/
package error

import (
	"fmt"
	"runtime"
)

// type ErrorCode represents the
type ErrorCode uint8

type Error struct {
	Code        ErrorCode
	Description string
	Stack       string
}

const (
	Configuration ErrorCode = iota
	Library
	MissingFile
	Sanity
	Value
	maxErrorCode

	defaultStackSize = 1 << 10
)

// codeToString returns a string representing the given error code.
func (code ErrorCode) String() string {
	switch code {
	case Configuration:
		return "Configuration Error"
	case Library:
		return "Library Error"
	case MissingFile:
		return "Missing File Error"
	case Sanity:
		return "Sanity Error"
	case Value:
		return "Value Error"
	}

	return fmt.Sprintf("Unrecognized Error Code %d", code)
}

// Error returns a string describing the given Error. It is not
// newline-terminaled.
func (err *Error) Error() string {
	if err == nil {
		return "Value Error: Error() called on nil pointer."
	}

	return fmt.Sprintf("%s: %s", err.Code.String(), err.Description)
}

// VerboseError returns a string containing both a stack trace and a string
// describing the error. It is not newline-terminated.
func (err *Error) VerboseError() string {
	if err == nil {
		return "Value Error: VerboseError() called on nil pointer."
	}

	return fmt.Sprintf("%s\n\n%s", err.Stack, err.Error())
}

// New creates a new Error corresponding to an error of type code which is
// described by the string desc.
func New(code ErrorCode, desc string) *Error {
	err := &Error{code, desc, ""}

	bytesRead, stackSize := defaultStackSize + 1, defaultStackSize
	var stackBuf []byte
	for stackSize < bytesRead {
		stackBuf = make([]byte, stackSize)
		bytesRead = runtime.Stack(stackBuf, false)
		stackSize = stackSize << 1
	}

	err.Stack = string(stackBuf[:bytesRead])
	return err
}

// Report prints an Error to stdout along with whatever other formatting is
// neccesary. This should only be used either as a last-ditch resort, like when
// setup of the GUI fails.
func Report(err *Error) {
	fmt.Println("A fatal error has occured.")
	fmt.Println()
	fmt.Println(err.VerboseError())
}
