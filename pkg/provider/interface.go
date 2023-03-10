package provider

import "os"

type CloudDisk interface {
	GetFile(path string) []FileInfo
}

type FileInfo interface {
	os.FileInfo
	Path() string
	Digest() string
}
