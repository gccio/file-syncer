package local

import (
	"os"
	"path"

	"github.com/gccio/file-syncer/pkg/provider"
	"go.uber.org/zap"
)

var _ provider.FileInfo = &localFile{}

type local struct {
	logger *zap.SugaredLogger
}

type localFile struct {
	os.FileInfo
	path   string
	digest string
}

func (l localFile) Digest() string {
	return l.digest
}
func (l localFile) Path() string {
	return l.path
}

func NewProvider(logger *zap.SugaredLogger) provider.CloudDisk {
	return &local{
		logger: logger,
	}
}

func (l *local) GetFile(p string) []provider.FileInfo {
	var localFileList []provider.FileInfo
	dirs, err := os.ReadDir(p)
	if err != nil {
		l.logger.Errorf("read local file %s failed with error %+v. skip this path.", p, err)
		return nil
	}

	for _, dir := range dirs {
		fp := path.Join(p, dir.Name())
		if dir.IsDir() {
			localFileList = append(localFileList, l.GetFile(fp)...)
			continue
		}
		fileInfo, err := os.Stat(fp)
		if err != nil {
			l.logger.Errorf("stat file %s failed with error %+v", fp, err)
			continue
		}
		// TODO 考虑何时计算md5值，若云盘上不存在，则不需要计算直接推送。
		localFileList = append(localFileList, &localFile{
			fileInfo,
			path.Join(p, fileInfo.Name()),
			"",
		})
	}

	l.logger.Infof("[end] read local file %s, count %d", p, len(localFileList))

	return localFileList
}

var _ provider.CloudDisk = &local{}
