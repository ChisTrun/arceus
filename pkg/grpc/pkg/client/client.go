package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	carbon "arceus/pkg/carbon/api"

	"google.golang.org/grpc/credentials"
	trace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const XRequestIDHeader = "x-request-id"

func NewClientConnection(config *carbon.TCPSocket, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	si := trace.StreamClientInterceptor(trace.WithServiceName(os.Getenv("DD_SERVICE")))
	ui := trace.UnaryClientInterceptor(trace.WithServiceName(os.Getenv("DD_SERVICE")))
	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(InjectRequestMetadata),
		grpc.WithStreamInterceptor(si),
		grpc.WithUnaryInterceptor(ui),
	}

	if config.GetSecure() {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	options = append(options, opts...)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", config.Address, config.GetPort()), options...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func InjectRequestMetadata(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	inMd, ok := metadata.FromIncomingContext(ctx)
	if ok {
		requestId := inMd.Get(XRequestIDHeader)
		md.Set(XRequestIDHeader, requestId...)
	}

	newCtx := metadata.NewOutgoingContext(ctx, md)
	return invoker(newCtx, method, req, reply, cc, opts...)
}
