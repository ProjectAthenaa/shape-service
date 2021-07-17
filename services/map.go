package services

import (
	"context"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	shape "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	shapegen "shape/shape"
)

var (
	trgt, _ = getGlobalHolder(context.Background(), shape.SITE_TARGET)
	rdb     = sonic.ConnectToRedis()
	Sites   = map[shape.SITE]string{
		shape.SITE_TARGET:    "shape:target",
		shape.SITE_END:       "shape:end",
		shape.SITE_NORDSTORM: "shape:nordstorm",
	}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func getGlobalHolder(ctx context.Context, site shape.SITE) (*shapegen.GlobalHolder, error) {
	val, err := rdb.Get(ctx, Sites[site]).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var gHolder *shapegen.GlobalHolder

	if err = json.Unmarshal([]byte(val), &gHolder); err != nil {
		return nil, err
	}

	return gHolder, nil
}
