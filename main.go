package main

import (
	"context"
	"time"

	"github.com/shivansh98/fampay-video-library/cron"
	"github.com/shivansh98/fampay-video-library/server"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go cron.CronRun(ctx, time.Minute, time.Now().Add(10*time.Minute).Unix(), "")
	defer cancel()
	server.Run()
}
