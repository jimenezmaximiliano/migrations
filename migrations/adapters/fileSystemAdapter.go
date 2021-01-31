package adapters

import (
	"io/ioutil"
	"os"
)

type FileSystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

type IOUtilAdapter struct{}

func (adapter IOUtilAdapter) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (adapter IOUtilAdapter) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
