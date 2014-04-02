package fpm

import (
	"io"
	"os"
)

type File interface {
	Name() string
	Data() (io.ReadCloser, error)
	Size() int64
	Mode() os.FileMode
}
