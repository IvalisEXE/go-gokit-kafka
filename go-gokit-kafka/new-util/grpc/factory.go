package grpc

import (
	"errors"
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

//ClientOption stores grpc client options
type ClientOption struct {
	//Timeout for circuit breaker
	Timeout time.Duration
	//Number of retry
	Retry int
	//Timeout for retry
	RetryTimeout time.Duration
}

func grpcConnection(address string, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	if creds == nil {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))
	}
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//EndpointFactory returns endpoint factory
func EndpointFactory(makeEndpoint func(*grpc.ClientConn, time.Duration, stdopentracing.Tracer, log.Logger) endpoint.Endpoint, creds credentials.TransportCredentials, timeout time.Duration, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		if instance == "" {
			return nil, nil, errors.New("Empty instance")
		}

		conn, err := grpcConnection(instance, creds)
		if err != nil {
			return nil, nil, err
		}
		endpoint := makeEndpoint(conn, timeout, tracer, logger)

		return endpoint, conn, nil
	}
}
