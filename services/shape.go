package services

import (
	"context"
	"github.com/ProjectAthenaa/shape"
	protos "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
)

type Server struct {
	protos.UnimplementedShapeServer
}

func (s Server) GenHeaders(ctx context.Context, site *protos.Site) (*protos.Headers, error) {
	if site.ResString != nil {
		return &protos.Headers{Values: shape.GenerateHeaders(site.Value, *site.ResString)}, nil
	}

	return &protos.Headers{Values: shape.GenerateHeaders(site.Value)}, nil
}
