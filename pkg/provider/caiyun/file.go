package caiyun

import (
	"fmt"
	"github.com/gccio/file-syncer/pkg/provider"
	"os"
	"path"
	"strings"
	"time"
)

const (
	cacheKeyFile = "file:"
	cacheKeyInfo = "info:"
	cachePeriod  = time.Hour
)

type fileLookupResult struct {
	file *FileInfo
	disk *DiskResult
	err  error
}

func (c *caiYun) cacheFileLookupResult(p string, lookup *fileLookupResult) {
	c.cache.Set(path.Join(cacheKeyFile, p), lookup, cachePeriod)
	disk := lookup.disk
	for _, catalog := range disk.CatalogList {
		c.cache.Set(path.Join(cacheKeyInfo, p, catalog.CatalogName), &FileInfo{
			id:      catalog.CatalogID,
			name:    catalog.CatalogName,
			path:    path.Join(p, catalog.CatalogName),
			isDir:   true,
			modTime: time.Time{},
			size:    0,
		}, cachePeriod)
	}
}

func (c *caiYun) GetFile(p string) []provider.FileInfo {
	f := c.getFile(p)
	if f.err != nil {
		c.logger.Error(f.err)
		return nil
	}

	files := make([]provider.FileInfo, 0)

	for _, content := range f.disk.ContentList {
		files = append(files, &FileInfo{
			id:      content.ContentID,
			name:    content.ContentName,
			path:    fmt.Sprintf("%s/%s", p, content.ContentName),
			isDir:   false,
			modTime: formatCaiYunDate(content.UpdateTime),
			size:    content.ContentSize,
			digest:  content.Digest,
		})
	}

	for _, catalog := range f.disk.CatalogList {
		files = append(files, c.GetFile(path.Join(p, catalog.CatalogName))...)
		//files = append(files, &FileInfo{
		//	id:      catalog.CatalogID,
		//	name:    catalog.CatalogName,
		//	path:    fmt.Sprintf("%s/%s", p, catalog.CatalogName),
		//	isDir:   true,
		//	modTime: formatCaiYunDate(catalog.UpdateTime),
		//	size:    0,
		//	digest:  "",
		//})
	}
	return files
}

func (c *caiYun) getFile(p string) *fileLookupResult {
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return c.getRootFile()
	}

	if lookup, found := c.cache.Get(path.Join(cacheKeyFile, p)); found {
		return lookup.(*fileLookupResult)
	}

	var (
		lookup   = &fileLookupResult{}
		fileInfo *FileInfo
	)

	obj, found := c.cache.Get(path.Join(cacheKeyInfo, p))
	if !found {
		r := c.getFile(path.Dir(p))
		if r.err != nil {
			return r
		}
		obj, found = c.cache.Get(path.Join(cacheKeyInfo, p))
		if !found {
			r.err = os.ErrNotExist
			return r
		}
	}
	fileInfo = obj.(*FileInfo)

	catalogID := fileInfo.id
	disk, err := c.GetDisk(catalogID)
	if err != nil {
		return &fileLookupResult{
			file: nil,
			disk: nil,
			err:  err,
		}
	}
	lookup.disk = disk
	lookup.file = fileInfo
	lookup.err = nil

	c.cacheFileLookupResult(p, lookup)
	return lookup
}
func (c *caiYun) getRootFile() *fileLookupResult {
	var (
		p    = "/"
		key  = cacheKeyFile + p
		name = "root"
	)

	if lookup, found := c.cache.Get(key); found {
		return lookup.(*fileLookupResult)
	}
	disk, err := c.GetDisk(defaultCatalog)
	if err != nil {
		return &fileLookupResult{
			file: nil,
			disk: nil,
			err:  err,
		}
	}
	fileInfo := &FileInfo{
		id:      defaultCatalog,
		name:    name,
		path:    p,
		isDir:   true,
		modTime: time.Time{},
		size:    0,
	}

	lookup := &fileLookupResult{
		file: fileInfo,
		disk: disk,
		err:  err,
	}
	c.cacheFileLookupResult("/", lookup)
	return lookup
}

type FileInfo struct {
	id      string
	name    string
	path    string
	isDir   bool
	modTime time.Time
	size    int64
	digest  string
}

var _ provider.FileInfo = &FileInfo{}

func (fi *FileInfo) Digest() string {
	return fi.digest
}

func (fi *FileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *FileInfo) Name() string {
	return fi.name
}

func (fi *FileInfo) Path() string {
	return fi.path
}

func (fi *FileInfo) Size() int64 {
	return fi.size
}

func (fi *FileInfo) Mode() os.FileMode {
	if fi.isDir {
		return os.ModeDir | 0o777
	}
	return 0o777
}

func (fi *FileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *FileInfo) Sys() interface{} {
	return fi
}

func (fi *FileInfo) ID() string {
	return fi.id
}
