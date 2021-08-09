package services

import (
	"context"
	shape "github.com/ProjectAthenaa/shape"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	protos "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
)

type Server struct {
	protos.UnimplementedShapeServer
}

func (s Server) GenHeaders(ctx context.Context, site *protos.Site) (*protos.Headers, error) {
	globalHolder, err := getGlobalHolder(ctx, siteConverter(site.Value))
	if err != nil {
		return nil, err
	}

	return &protos.Headers{Values: shape.GenerateHeaders(globalHolder)}, nil
}

func siteConverter(site protos.SITE) string {
	switch site {
	case protos.SITE_TARGET:
		return sonic.TARGET
	case protos.SITE_NEWBALANCE:
		return sonic.NEWBALANCE
	case protos.SITE_END:
		return sonic.END
	}
	return ""
}
