package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB connects with the MongoDB atlas and returns collection
func ConnectDB(ctx context.Context) (colln *mongo.Collection, err error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://shivansh:LtAAZn1fjWbiAUKT@cluster0.nhxu12q.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	colln = client.Database("fampay-db").Collection("_video_store")
	return colln, nil
}
