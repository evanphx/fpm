package fpm

import (
	"io"
	"os"
	"path/filepath"
)

type DirPackage struct {
	Path string
}

type DirPackageFile struct {
	path, name string
	info       os.FileInfo
}

func (d *DirPackageFile) Name() string {
	return d.name
}

func (d *DirPackageFile) Data() (io.ReadCloser, error) {
	return os.Open(d.path)
}

func (d *DirPackageFile) Mode() os.FileMode {
	return d.info.Mode()
}

func (d *DirPackageFile) Size() int64 {
	return d.info.Size()
}

func (d *DirPackage) Files() []File {
	var files []File

	filepath.Walk(d.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(d.Path, path)
		if err != nil {
			return err
		}

		f := &DirPackageFile{path, rel, fi}
		files = append(files, f)

		return nil
	})

	return files
}
