package services

import (
	"context"
	"encoding/json"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/go-redis/redis/v8"
	shape "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	shapegen "shape/shape"
)

var (
	rdb   = sonic.ConnectToRedis()
	Sites = map[shape.SITE]string{
		shape.SITE_TARGET:    "shape:target",
		shape.SITE_END:       "shape:end",
		shape.SITE_NORDSTORM: "shape:nordstorm",
	}
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
