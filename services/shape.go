package services

import (
	"context"
	shape "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"log"
	"time"
)

type Server struct {
	shape.UnimplementedShapeServer
}

func (s Server) GenHeaders(ctx context.Context, site *shape.Site) (*shape.Headers, error) {
	start := time.Now()
	//globalHolder, err := getGlobalHolder(ctx, site.Value)
	hdrs := trgt.GenerateHeaders()
	log.Printf("Generated Headers | %s", time.Since(start))
	//if err != nil {
	//	log.Println(err)
	//	return nil, err
	//}

	return &shape.Headers{Values: hdrs}, nil
}
