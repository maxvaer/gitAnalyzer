// Package Utils provides a collection of utility functions for various tasks.
package Utils

import (
	"errors"
	"io"
	"os"
	"strings"
)

// Contains checks if a given slice slc contains a given string str.
func Contains(slc []string, str string) bool {
	for _, value := range slc {
		if value == str {
			return true
		}
	}
	return false
}

// FileExists checks if a given file path exists on the system.
func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// NormalizeGitURLToHTTPS converts a given GitHub SSH URL to the https scheme.
func NormalizeGitURLToHTTPS(url string) string {
	url = strings.Replace(url, "git://", "https://", 1)
	url = strings.Replace(url, "git@", "https://", 1)
	// Count the number of :
	colonCount := strings.Count(url, ":")
	if colonCount == 2 {
		// Given url has the format https://github.com:maxvaer/scanner-test-repo
		// therefore replace the : after the username
		lastColonIndex := strings.LastIndex(url, ":")
		url = url[:lastColonIndex] + strings.Replace(url[lastColonIndex:], ":", "/", 1)
	}
	return url
}

// IsDirEmpty checks if a given directory path is an existing directory.
func IsDirEmpty(path string) (bool, error) {
	// try to open the given path
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	// try to read content of the directory
	_, err = dir.Readdir(1)
	if err == io.EOF {
		// if EOF, directory has no content, therefore is empty
		return true, nil
	}
	// no EOF, directory is not empty
	return false, err
}
