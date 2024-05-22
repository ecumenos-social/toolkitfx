package fxgrpc

import (
	"context"
	"net"
	"net/http"
	"time"

	grpcutils "github.com/ecumenos-social/grpc-utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type HTTPGatewayConfig struct {
	Host string `default:"0.0.0.0"`
	Port string `default:"9090"`
}

type HTTPGatewayHandler struct {
	Handler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

func RunHTTPGateway(lc fx.Lifecycle, s fx.Shutdowner, logger *zap.Logger, config *Config, g *HTTPGatewayHandler) error {
	addr := net.JoinHostPort(config.HTTPGateway.Host, config.HTTPGateway.Port)
	mux := runtime.NewServeMux()
	conn := grpcutils.NewClientConnection(net.JoinHostPort(config.GRPC.Host, config.GRPC.Port))

	var httpServer *http.Server
	localCtx := context.Background()
	_ = conn.Dial(grpcutils.DefaultDialOpts(logger)...)
	err := g.Handler(localCtx, mux, conn.Connection)
	if err != nil {
		logger.Error("Failed to register mapping service handler", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("starting HTTP gateway", zap.String("addr", addr))
				httpServer = &http.Server{Addr: addr, Handler: mux}
				err = httpServer.ListenAndServe()
				if err != nil {
					logger.Error("Failed to run HTTP gateway", zap.Error(err))
					_ = s.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ = conn.CleanUp()
			if httpServer != nil {
				timeout, can := context.WithTimeout(context.Background(), 10*time.Second)
				defer can()
				if err := httpServer.Shutdown(timeout); err != nil {
					logger.Error("Stopped http server after grpc failure", zap.Error(err))
				}
			}
			return nil
		},
	})

	return nil
}
