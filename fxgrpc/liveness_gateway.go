package fxgrpc

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LivenessGatewayConfig struct {
	Host string `default:"0.0.0.0"`
	Port string `default:"8086"`
}

type LivenessGatewayHandler struct {
	Handler http.Handler
}

func RunLivenessGateway(lc fx.Lifecycle, s fx.Shutdowner, logger *zap.Logger, config *Config, h *LivenessGatewayHandler) error {
	addr := net.JoinHostPort(config.LivenessGateway.Host, config.LivenessGateway.Port)
	var httpServer *http.Server

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("starting liveness gateway...", zap.String("addr", addr))
				httpServer = &http.Server{Addr: addr, Handler: h.Handler}
				err := httpServer.ListenAndServe()
				if err != nil {
					logger.Error("Failed to run liveness gateway", zap.Error(err))
					_ = s.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
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
