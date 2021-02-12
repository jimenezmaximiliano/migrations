package adapters

import (
	"io/ioutil"
	"os"
)

// FileSystem is an interface to read contents from the file system.
type FileSystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

// IOUtilAdaper is an implementation of FileSystem using io/ioutil and os.
type IOUtilAdapter struct{}

// Ensure IOUtilAdapter implements FileSystem
var _ FileSystem = IOUtilAdapter{}

// ReadDir list the files on a given directory.
func (adapter IOUtilAdapter) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// ReadFile reads the contents of a file.
func (adapter IOUtilAdapter) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
