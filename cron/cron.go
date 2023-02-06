package cron

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shivansh98/fampay-video-library/database"
	"github.com/shivansh98/fampay-video-library/util"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const (
	developerKey = "AIzaSyC9VYkoq4KLY-jf9-7v_2NKB5pxnK0F7_8"
	query        = "Cricket"
	maxResults   = 25
)

// CronRun runs a recursive function till the time is allowed (initially 10 min with delay of 1 min each) , it recursively
// fetches the next page of the youtube api results
func CronRun(ctx context.Context, delay time.Duration, nextPageToken string) {
	select {
	case <-ctx.Done():
		log.Print("ctx is cancelled , closing cron")
		return
	default:
	}
	nextPageToken = fetchYoutubeVideos(ctx, nextPageToken)
	if nextPageToken != "" {
		time.Sleep(delay)
		CronRun(ctx, delay, nextPageToken)
		return
	}
}

// fetchYoutubeVideos fetches data from youtube data api by appyling the next page token
func fetchYoutubeVideos(ctx context.Context, getPage string) (nextPageToken string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		util.LogError(err, "creating new api client")
		return getPage
	}

	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults).Type("video").PageToken(getPage)
	response, err := call.Do()
	if err != nil {
		util.LogError(err, "parsing response")
		return getPage
	}

	dbReq, err := database.CreateDBRequest(response)
	if err != nil {
		util.LogError(err, "making db request")
		return getPage
	}

	err = dbReq.Insert(ctx)
	if err != nil {
		util.LogError(err, "inserting data into db")
		return getPage
	}
	fmt.Println("Cron Iteration ran , Timestamp ", time.Now(), " : inteserted ", len(response.Items), " items in DB")
	nextPageToken = response.NextPageToken
	return
}
