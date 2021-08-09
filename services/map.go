package services

import (
	"context"
	"github.com/ProjectAthenaa/shape/deobfuscation"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/frame"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"os"
	"shape/debug"
)

var (
	trgt, _ = getGlobalHolder(context.Background(), sonic.TARGET)
	rdb     = frame.ConnectRedis(os.Getenv("REDIS_URL"))
	Sites   = map[string]string{
		sonic.TARGET:     "shape:target",
		sonic.END:        "shape:end",
		sonic.NORDSTORM:  "shape:nordstorm",
		sonic.NEWBALANCE: "shape:newbalance",
	}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func getGlobalHolder(ctx context.Context, site string) (*deobfuscation.GlobalHolder, error) {
	if os.Getenv("DEBUG") == "1" {
		switch site {
		case sonic.TARGET:
			return debug.Target.GlobalHolder, nil
		case sonic.NEWBALANCE:
			return debug.NewBalance.GlobalHolder, nil
		}
	}

	val, err := rdb.Get(ctx, Sites[site]).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var gHolder *deobfuscation.GlobalHolder

	if err = json.Unmarshal([]byte(val), &gHolder); err != nil {
		return nil, err
	}

	return gHolder, nil
}
