package server

import (
	"context"
	"encoding/json"
	"fampay-video-library/database"
	"fampay-video-library/util"
	"fmt"
	"io"
	"net/http"
)

func SearchmiddleWare(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		searchVideos(w, r)
		return
	}
	failure(fmt.Errorf("invalid http method used"), w)
}

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
	w.Write(resp)
}

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
