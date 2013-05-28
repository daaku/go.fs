// Package fsutil provides utilities for working with File Systems.
package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var errInvalidCharacterInPath = errors.New("invalid character in file path")

type errNotFound string

func (e errNotFound) Error() string {
	return fmt.Sprintf("file not found: %s", string(e))
}

// Returns an error that indicates the named file was not found.
func NewErrNotFound(name string) error {
	return errNotFound(name)
}

type errLimitedNotFound string

func (e errLimitedNotFound) Error() string {
	return fmt.Sprintf("file limited or not found: %s", string(e))
}

// Returns an error that indicates the named file was not found or access was
// artificially limited.
func NewErrLimitedNotFound(name string) error {
	return errLimitedNotFound(name)
}

func IsNotExist(err error) bool {
	if _, ok := err.(errNotFound); ok {
		return true
	}
	if _, ok := err.(errLimitedNotFound); ok {
		return true
	}
	return os.IsNotExist(err)
}

// Cleans path string.
func Clean(name string) (string, error) {
	if filepath.Separator != '/' &&
		strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return "", errInvalidCharacterInPath
	}
	return filepath.FromSlash(path.Clean("/" + name)), nil
}
