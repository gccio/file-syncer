package sync

import (
	"github.com/gccio/file-syncer/pkg/provider"
	"github.com/gccio/file-syncer/pkg/provider/caiyun"
	"github.com/gccio/file-syncer/pkg/provider/local"
	"time"

	"github.com/pmylund/go-cache"
)

var _ FileManager = &fileManager{}

type FileManager interface {
	SrcFile(string) []provider.FileInfo
	DestFile(string) []provider.FileInfo
}

type fileManager struct {
	cache *cache.Cache
	src   provider.CloudDisk
	dest  provider.CloudDisk
}

func NewFileManager() FileManager {
	return &fileManager{
		cache: cache.New(time.Minute*10, time.Minute),
		src:   local.NewProvider(sugar),
		dest:  caiyun.NewProvider(sugar),
	}
}

func (f fileManager) SrcFile(p string) []provider.FileInfo {
	return f.src.GetFile(p)
}

func (f fileManager) DestFile(p string) []provider.FileInfo {
	return f.dest.GetFile(p)
}
