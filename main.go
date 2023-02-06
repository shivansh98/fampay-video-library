package main

import (
	"context"
	"time"

	"github.com/shivansh98/fampay-video-library/cron"
	"github.com/shivansh98/fampay-video-library/server"
)

func main() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Minute*10)
	go cron.CronRun(ctx, time.Minute, "")
	server.Run()
}
