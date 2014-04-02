package fpm

type Package interface {
	Files() []File
	Write(input Package) error
}
