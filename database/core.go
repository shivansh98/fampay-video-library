package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/youtube/v3"
)

type DBRequest struct {
	Items []*ItemDetails `bson:"items"`
}

type ItemDetails struct {
	Description string           `bson:"description" json:"description"`
	PublishedAt string           `bson:"published_at" json:"published_at"`
	Thumbnails  ThumbnailDetails `bson:"thumbnails" json:"thumbnails"`
	Title       string           `bson:"title" json:"title"`
}

type ThumbnailDetails struct {
	Default  Url `bson:"default" json:"default"`
	High     Url `bson:"high" json:"high"`
	Maxres   Url `bson:"max_res" json:"max_res" `
	Medium   Url `bson:"medium" json:"medium"`
	Standard Url `bson:"standard" json:"standard"`
}

type Url string

const limit = 25

// CreateDBRequest creates a DB request from the youtube data api response, to insert data into mongoDB
func CreateDBRequest(apiResp *youtube.SearchListResponse) (*DBRequest, error) {
	if apiResp == nil {
		return nil, fmt.Errorf("got api response as nil")
	}
	req := DBRequest{}
	req.Items = make([]*ItemDetails, len(apiResp.Items))
	for _, item := range apiResp.Items {
		if item == nil || item.Snippet == nil || item.Snippet.Thumbnails == nil {
			continue
		}

		i := ItemDetails{
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
			Title:       item.Snippet.Title,
		}
		if item.Snippet.Thumbnails.Default != nil {
			i.Thumbnails.Default = Url(item.Snippet.Thumbnails.Default.Url)
		}
		if item.Snippet.Thumbnails.High != nil {
			i.Thumbnails.High = Url(item.Snippet.Thumbnails.High.Url)
		}
		if item.Snippet.Thumbnails.Maxres != nil {
			i.Thumbnails.Maxres = Url(item.Snippet.Thumbnails.Maxres.Url)
		}
		if item.Snippet.Thumbnails.Medium != nil {
			i.Thumbnails.Medium = Url(item.Snippet.Thumbnails.Medium.Url)
		}
		if item.Snippet.Thumbnails.Standard != nil {
			i.Thumbnails.Standard = Url(item.Snippet.Thumbnails.Standard.Url)
		}
		req.Items = append(req.Items, &i)
	}
	return &req, nil
}

// Insert inserts the data got from youtube api
func (d *DBRequest) Insert(ctx context.Context) error {
	colln, err := ConnectDB(ctx)
	if err != nil {
		return err
	}
	dc := make([]interface{}, 0)

	for _, item := range d.Items {
		if item == nil {
			continue
		}
		dc = append(dc, interface{}(*item))
	}
	if len(dc) == 0 {
		return fmt.Errorf("no doc present in DB")
	}
	_, err = colln.InsertMany(ctx, dc)
	return err
}

// GetDocuments calls fetchDocs with appropriate params to fetch docs
func GetDocuments(ctx context.Context, pageNo int64) (items []*ItemDetails, err error) {
	return fetchDocs(ctx, "", pageNo, "/get-videos")
}

// SearchDocuments calls fetchDocs with appropriate params to fetch docs
func SearchDocuments(ctx context.Context, text string) (items []*ItemDetails, err error) {
	return fetchDocs(ctx, text, 0, "/search")

}

// fetchDocs fetches doc from DB based on the http request params
// for /search -> it uses text param only
// searches doucments based on the text supplied, it used the text indexes
// of the mongoDB to search . text index in mongoDB is formed by combining title and description field together.

// for /get-videos -> it uses pageNo param only
// here it skips the documents
// that come before that page_no and sort the documents based on published_at in descending order.
func fetchDocs(ctx context.Context, text string, pageNo int64, request string) (items []*ItemDetails, err error) {
	colln, err := ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	opts := options.Find().SetSort(bson.D{{Key: "published_at", Value: -1}}).SetLimit(limit)

	var query interface{}
	switch request {
	case "/search":
		query = bson.M{
			"$text": bson.M{
				"$search": text,
			},
		}
	case "/get-videos":
		query = bson.D{}
		if pageNo != 0 {
			opts.SetSkip(limit * pageNo)
		}
	}

	cur, err := colln.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	items = []*ItemDetails{}
	err = cur.All(ctx, &items)
	if err != nil {
		return nil, err
	}
	return
}
