// Package fsutil provides utilities for working with File Systems.
package fsutil

import (
	"errors"
	"fmt"
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

func IsNotExist(err error) bool {
	_, ok := err.(errNotFound)
	return ok
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
