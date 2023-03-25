package mongodb

import (
	"context"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Client *mongo.Client
var DB *mongo.Database

func Setup() {
	if config.Get("mongodb") == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.GetString("mongodb.url")))
	if err != nil {
		logger.Panicf("fail to connect mongodb: %s", err)
	}
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()
	err = Client.Ping(pingCtx, nil)
	if err != nil {
		logger.Fatal(err)
	}
	DB = Client.Database(config.GetString("mongodb.database"))
}
