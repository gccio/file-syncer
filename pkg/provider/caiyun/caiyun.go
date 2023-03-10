package caiyun

import (
	"github.com/gccio/file-syncer/pkg/provider"
	"github.com/pmylund/go-cache"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type caiYun struct {
	logger *zap.SugaredLogger
	cache  *cache.Cache
	cli    *http.Client
}

var _ provider.CloudDisk = &caiYun{}

func NewProvider(logger *zap.SugaredLogger) provider.CloudDisk {
	return &caiYun{
		logger: logger,
		cache:  cache.New(time.Minute*10, time.Minute),
		cli: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0 * time.Second,
		},
	}
}
