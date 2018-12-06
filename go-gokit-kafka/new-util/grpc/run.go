package grpc

import (
	"errors"
	"net"
	"net/http"
	"strings"

	util "git.bluebird.id/bluebird/util/log"
	"github.com/go-kit/kit/log"

	middle "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	//ErrServer is internal server error
	ErrServer = errors.New("Internal server error")
)

//Recovery return grpc server option with recovery handler
func Recovery(logger log.Logger) []grpc.ServerOption {
	handler := func(p interface{}) (err error) {
		logger.Log("panic", p)
		return ErrServer
	}
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(handler),
	}
	serverOptions := []grpc.ServerOption{
		middle.WithUnaryServerChain(
			recovery.UnaryServerInterceptor(opts...),
		),
		middle.WithStreamServerChain(
			recovery.StreamServerInterceptor(opts...),
		)}
	return serverOptions
}

//Serve listen for client request
func Serve(address string, server *grpc.Server, logger log.Logger) {

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Log(util.LogError, err.Error())
		return
	}

	err = server.Serve(lis)
	if err != nil {
		logger.Log(util.LogError, err.Error())
		return
	}
}

//RegisterHTTPHandler register endpoint to http server
type RegisterHTTPHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

//ServeHTTP listen for http request
func ServeHTTP(grpcAddress, httpAddress string, register RegisterHTTPHandler,
	creds credentials.TransportCredentials, logger log.Logger, allowcors bool) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	var opts []grpc.DialOption
	if creds == nil {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	} else {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	}
	err := register(ctx, mux, grpcAddress, opts)
	if err != nil {
		logger.Log(util.LogError, err.Error())
		return
	}

	if allowcors {
		http.ListenAndServe(httpAddress, allowCORS(mux))
	} else {
		http.ListenAndServe(httpAddress, mux)
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept", "Authorization"}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
