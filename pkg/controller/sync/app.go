package sync

import (
	"context"
	"time"

	"github.com/gccio/file-syncer/pkg/provider"
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
)

type Sync interface {
	Run()
}

type syncDir struct {
	Src  string `json:"local,omitempty"`
	Dest string `json:"remote,omitempty"`
}

type sync struct {
	ctx context.Context

	logger      *zap.SugaredLogger
	syncDir     []syncDir
	fileManager FileManager
	deleteCh    []chan int
	uploadCh    []chan int
}

func NewSyncController(ctx context.Context) Sync {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar = logger.Sugar()

	return &sync{
		ctx:         ctx,
		logger:      sugar,
		syncDir:     []syncDir{},
		fileManager: NewFileManager(),
	}
}

func (s *sync) Run() {
	s.process()

	tk := time.NewTimer(time.Minute * 10)
	go s.uploadWorker()
	go s.deleteWorker()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tk.C:
			s.process()
		}
	}
}

func (s *sync) process() {
	for _, dir := range s.syncDir {
		sugar.Infof("%+v", dir)
		srcFile := s.fileManager.SrcFile(dir.Src)
		destFile := s.fileManager.DestFile(dir.Dest)
		s.logger.Infof("src: %+v", srcFile)
		s.logger.Infof("dest %+v", destFile)
	}
}

type Operation int

func (s *sync) merge(srcFile, destFile []provider.FileInfo) {
	srcMap := map[string]provider.FileInfo{}
	for _, src := range srcFile {
		srcMap[src.Path()] = src
	}
	for _, dest := range destFile {
		// 数据在本地已被删除，云端同样应该被删除
		if _, ok := srcMap[dest.Path()]; !ok {

		}
	}
}

func (s *sync) uploadWorker() {

}

func (s *sync) deleteWorker() {

}
