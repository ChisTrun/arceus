package client

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	trace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type GrpcClientConn struct {
	conn         *grpc.ClientConn
	methodPrefix string
}

// type check
var _ grpc.ClientConnInterface = (*GrpcClientConn)(nil)

// Invoke implements grpc.ClientConnInterface.
func (c *GrpcClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	newMethod := path.Join(c.methodPrefix, method)
	return c.conn.Invoke(ctx, newMethod, args, reply, opts...)
}

// NewStream implements grpc.ClientConnInterface.
func (c *GrpcClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	newMethod := path.Join(c.methodPrefix, method)
	return c.conn.NewStream(ctx, desc, newMethod, opts...)
}

func NewClientConnectionV2(target string, opts ...grpc.DialOption) (*GrpcClientConn, error) {
	parsedUrl, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	scheme := parsedUrl.Scheme
	path := parsedUrl.Path
	port := parsedUrl.Port()
	if port == "" {
		if scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	address := fmt.Sprintf("%s:%s", parsedUrl.Hostname(), port)

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

	if scheme == "https" {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	options = append(options, opts...)

	conn, err := grpc.DialContext(ctx, address, options...)
	if err != nil {
		return nil, err
	}
	return &GrpcClientConn{
		conn:         conn,
		methodPrefix: path}, nil
}
