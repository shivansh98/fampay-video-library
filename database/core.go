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

func GetDocuments(ctx context.Context, pageNo int64) (items []*ItemDetails, err error) {
	colln, err := ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	opts := options.Find().SetSort(bson.D{{"published_at", -1}}).SetLimit(limit)
	if pageNo != 0 {
		opts.SetSkip(limit * pageNo)
	}
	cur, err := colln.Find(ctx, bson.D{}, opts)
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

func SearchDocuments(ctx context.Context, text string) (items []*ItemDetails, err error) {
	colln, err := ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	opts := options.Find().SetSort(bson.D{{Key: "published_at", Value: -1}}).SetLimit(limit)
	query := bson.M{
		"$text": bson.M{
			"$search": text,
		},
	}
	cur, err := colln.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	fmt.Println("count : ", cur.RemainingBatchLength())
	items = []*ItemDetails{}
	err = cur.All(ctx, &items)
	if err != nil {
		return nil, err
	}
	return
}
