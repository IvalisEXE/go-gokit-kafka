package opentracing

import (
	"git.bluebird.id/bluebird/tracer-go"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/credentials"
)

//Tracer returns default tracer
func Tracer(address, service string, tls credentials.TransportCredentials) opentracing.Tracer {
	return tracer.NewTracer(address, service, tls)
}
