package main

import (
	"github.com/evanphx/fpm/go/fpm"
	"os"
)

func main() {
	md := fpm.NewMetaData()

	dir := &fpm.DirPackage{os.Args[1]}
	deb := &fpm.DebPackage{md, os.Args[2]}

	err := deb.Write(dir)
	if err != nil {
		panic(err)
	}

}
