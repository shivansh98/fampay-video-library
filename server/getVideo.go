package server

import (
	"context"
	"encoding/json"

	"fmt"
	"net/http"
	"strconv"

	"github.com/shivansh98/fampay-video-library/database"
)

func getVideoMiddleWare(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getVideos(w, r)
		return
	}
	failure(fmt.Errorf("invalid http method used"), w)
}

func getVideos(w http.ResponseWriter, r *http.Request) {
	pageToken := r.URL.Query().Get(page_token)
	page_no := 0
	var err error
	if pageToken != "" {
		page_no, err = strconv.Atoi(pageToken)
		if err != nil {
			failure(fmt.Errorf("page token not an integer"), w)
			return
		}
	}
	ctx := context.Background()
	items, err := database.GetDocuments(ctx, int64(page_no))
	if err != nil {
		failure(err, w)
		return
	}

	if items == nil {
		failure(fmt.Errorf("got empty result"), w)
		return
	}
	httpResp := Response{
		Status: "success",
		Result: items,
		Page_Info: &PageInfo{
			Next: getPageURL(page_no + 1),
			Prev: getPageURL(page_no - 1),
		},
	}
	writeResp, err := json.Marshal(httpResp)
	if err != nil {
		failure(err, w)
		return
	}
	w.Write(writeResp)
	w.Header().Set("content-type", "application/json")
}
