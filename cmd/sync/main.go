package main

import (
	"context"
	"flag"
	"github.com/gccio/file-syncer/pkg/controller/sync"
)

func main() {
	flag.Parse()

	ctx := context.TODO()
	sync.NewSyncController(ctx).Run()
}
