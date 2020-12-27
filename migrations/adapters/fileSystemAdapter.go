package adapters

import (
	"io/ioutil"
	"os"
)

type FileSystemAdapter struct{}

func (adapter FileSystemAdapter) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (adapter FileSystemAdapter) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
