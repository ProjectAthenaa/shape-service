package services

import (
	"context"
	shape "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
)

type Server struct {
	shape.UnimplementedShapeServer
}

func (s Server) GenHeaders(ctx context.Context, site *shape.Site) (*shape.Headers, error) {
	globalHolder, err := getGlobalHolder(ctx, site.Value)
	if err != nil {
		return nil, err
	}

	return &shape.Headers{Values: globalHolder.GenerateHeaders()}, nil
}
