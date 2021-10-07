package services

import (
	"context"
	"errors"
	"github.com/ProjectAthenaa/shape"
	protos "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
)

type Server struct {
	protos.UnimplementedShapeServer
}

func (s Server) GenHeaders(ctx context.Context, site *protos.Site) (headers *protos.Headers, err error) {
	err = nil
	defer func() {
		if a := recover(); a != nil {
			err = errors.New("internal_error")
		}
	}()
	if site.ResString != nil {
		return &protos.Headers{Values: shape.GenerateHeaders(site.Value, *site.ResString)}, err
	}

	return &protos.Headers{Values: shape.GenerateHeaders(site.Value)}, err
}
