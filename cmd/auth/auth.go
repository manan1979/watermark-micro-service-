package main

import (
	//	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/manan1979/watermark-service/pkg/auth"
	"github.com/manan1979/watermark-service/pkg/auth/endpoint"
	"github.com/manan1979/watermark-service/pkg/auth/transport"
	"github.com/oklog/oklog/pkg/group"
)

const (
	defaultHTTPPort = "8081"
	defaultGRPCPort = "8082"
)

func main() {

	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHTTPPort))
	)

	var (
		authService     = auth.NewService()
		authEndpoints   = endpoint.NewAuthEndpointSet(authService)
		authHTTPHandler = transport.NewHTTPHandler(authEndpoints)
	)

	var g group.Group

	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, authHTTPHandler)
		}, func(error) {
			httpListener.Close()
		})
	}

	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				logger.Log("received signal", sig)
			case <-cancelInterrupt:
			}
			return nil
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
