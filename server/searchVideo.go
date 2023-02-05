package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shivansh98/fampay-video-library/database"
	"github.com/shivansh98/fampay-video-library/util"
)

// SearchMiddleWare is the middleware for /search request and checks if the request is POST then forward to
// main halderFunc searchVideos  otherwise returns error response
func SearchMiddleWare(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Add("Content-Type", "application/json")
		searchVideos(w, r)
		return
	}
	failure(fmt.Errorf("invalid http method used"), w)
}

// searchVideos is the main handlerFunc that handles the /search request and returns the requested response
// sorted by published_at in descending order
func searchVideos(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.LogError(fmt.Errorf("got error in getting body error:", err), "/search req")
		failure(fmt.Errorf("got error in getting body error:", err), w)
		return
	}
	var srchReq SearchRequest
	err = json.Unmarshal(body, &srchReq)
	if err != nil {
		util.LogError(fmt.Errorf("got error in unmarshalling body error:", err), "/search req")
		failure(fmt.Errorf("got error in unmarshalling body error:", err), w)
		return
	}

	items, err := database.SearchDocuments(context.Background(), srchReq.Text)
	if err != nil {
		util.LogError(fmt.Errorf("got error in getting documents from db error:", err), "/search req")
		failure(fmt.Errorf("got error in getting documents from db error:", err), w)
		return
	}
	resp, err := createSearchResp(items)
	if err != nil {
		util.LogError(fmt.Errorf("got error in creating search response error:", err), "/search req")
		failure(fmt.Errorf("got error in creating search response error:", err), w)
		return
	}
	fmt.Println("Got a request in /search , no. of results returning : ", len(items))
	w.Write(resp)
}

// createSearchResp creates a search Response for request
func createSearchResp(items []*database.ItemDetails) (resp []byte, err error) {
	if items == nil {
		return nil, err
	}
	httpResp := Response{
		Status: "success",
		Result: items,
	}
	resp, err = json.Marshal(httpResp)
	if err != nil {
		return nil, err
	}
	return
}
