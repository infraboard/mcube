package file

import (
	"os"
)

func stat(name string, statFunc func(name string) (os.FileInfo, error)) (Info, error) {
	info, err := statFunc(name)
	if err != nil {
		return nil, err
	}

	return wrap(info)
}

func wrap(info os.FileInfo) (Info, error) {
	return fileInfo{FileInfo: info}, nil
}
