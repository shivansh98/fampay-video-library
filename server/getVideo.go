package server

import (
	"context"
	"encoding/json"
	"log"

	"fmt"
	"net/http"
	"strconv"

	"github.com/shivansh98/fampay-video-library/database"
	"github.com/shivansh98/fampay-video-library/util"
)

// getVideoMiddleWare checks if the method of the request is other than GET then returns error response otherwise
// forward to the main handlerFunc getvideos
func getVideoMiddleWare(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", "application/json")
		getVideos(w, r)
		return
	}
	failure(fmt.Errorf("invalid http method used"), w)
}

// getVideos is the handlerFunc that serves the /get-videos request and returns a paginated response
func getVideos(w http.ResponseWriter, r *http.Request) {
	pageToken := r.URL.Query().Get(page_token)
	page_no := 0
	var err error
	if pageToken != "" {
		page_no, err = strconv.Atoi(pageToken)
		if err != nil {
			util.LogError(fmt.Errorf("page token not an integer"), "/get-videos req")
			failure(fmt.Errorf("page token not an integer"), w)
			return
		}
	}
	ctx := context.Background()
	items, err := database.GetDocuments(ctx, int64(page_no))
	if err != nil {
		util.LogError(err, "/get-videos req")
		failure(err, w)
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
		util.LogError(err, "/get-videos req")
		failure(err, w)
		return
	}
	log.Println("Got a request in /get-videos , no. of results returning : ", len(items))
	w.Write(writeResp)
}
