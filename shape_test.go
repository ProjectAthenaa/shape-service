package main_test

import (
	"context"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"shape/services"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	shape.RegisterShapeServer(server, services.Server{})
	go func() {
		server.Serve(lis)
	}()
}

func TestShape(t *testing.T)  {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil{
		t.Fatal(err)
	}


	client := shape.NewShapeClient(conn)

	headers, err := client.GenHeaders(context.Background(), &shape.Site{Value: shape.SITE_NEWBALANCE})
	if err != nil{
		t.Fatal(err)
	}


	fmt.Println(headers.Values)

}