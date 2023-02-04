package main

import (
	"context"
	"fampay-video-library/cron"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go cron.CronRun(ctx, time.Minute, time.Now().Add(10*time.Minute).Unix(), "")
	defer cancel()
}
