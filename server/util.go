package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shivansh98/fampay-video-library/database"
)

const page_token = "page_no"
const servURL = "http://localhost:8080/get-videos?page_no="

type Response struct {
	Status    string                  `json:"status,omitempty"`
	Message   string                  `json:"message,omitempty"`
	Result    []*database.ItemDetails `json:"result,omitempty"`
	Page_Info *PageInfo               `json:"page_info,omitempty"`
}

type SearchRequest struct {
	Text string `json:"text"`
}

type PageInfo struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

func failure(err error, w http.ResponseWriter) {
	var resp Response
	resp.Message = err.Error()
	resp.Status = "failure"
	b, _ := json.Marshal(resp)
	w.Write(b)
}

func getPageURL(page int) string {
	if page < 0 {
		return ""
	}
	return servURL + strconv.Itoa(page)
}
