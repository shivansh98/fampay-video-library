package cron

import (
	"context"
	"encoding/json"
	"fampay-video-library/database"
	"fampay-video-library/util"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const (
	developerKey = "AIzaSyC9VYkoq4KLY-jf9-7v_2NKB5pxnK0F7_8"
	query        = "Cricket"
	maxResults   = 25
)

func CronRun(ctx context.Context, delay time.Duration, till int64, nextPageToken string) error {
	if time.Now().Unix() >= till {
		return fmt.Errorf("done")
	}
	nextPageToken = fetchYoutubeVideos(ctx, nextPageToken)
	if nextPageToken != "" {
		time.Sleep(delay)
		return CronRun(ctx, delay, till, nextPageToken)
	}
	return fmt.Errorf("error in getting next page")
}

func fetchYoutubeVideos(ctx context.Context, getPage string) (nextPageToken string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	util.HandleError(err, "got error")

	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults).Type("video").PageToken(getPage)
	response, err := call.Do()
	util.HandleError(err, "got error")
	js, err := json.Marshal(response)
	util.HandleError(err, "got error")
	fmt.Print("REsponse from youtube api ", string(js))
	dbReq, err := database.CreateDBRequest(response)
	util.HandleError(err, "got error")

	err = dbReq.Insert(ctx)
	util.HandleError(err, "got error in inserting document in db")

	nextPageToken = response.NextPageToken
	return
}
