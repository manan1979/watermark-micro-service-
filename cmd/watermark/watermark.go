package watermark

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"

	pb "github.com/manan1979/watermark-service/api/pb/watermark"
	"github.com/manan1979/watermark-service/pkg/watermark"
	"github.com/manan1979/watermark-service/pkg/watermark/endpoint"
	"github.com/manan1979/watermark-service/pkg/watermark/transport"
	"google.golang.org/grpc"
)

const (
	defaultHTTPPort = "8081"
	defaultGRPCPort = "8082"
)

func main() {

	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("locahost", envString("HTTP_PORT", defaultHTTPPort))
		grpcAddr = net.JoinHostPort("localhost", envString("GRPC_PORT", defaultGRPCPort))
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC())

	var (
		service     = watermark.NewService()
		eps         = endpoint.NewEndpointSet(service)
		httpHandler = transport.NewHTTPHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
	)

	var g group.Group

	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", "httpAddr")
			return http.Serve(httpListener, httpHandler)

		}, func(error) {
			httpListener.Close()
		})
	}

	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", "grpcAddr")

			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterWatermarkServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}

	{
		cancelInterrupt := make(chan struct{})

		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())

}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
