package fxgrpc

import (
	"context"
	"net"

	grpcutils "github.com/ecumenos-social/grpc-utils"
	"github.com/ecumenos-social/toolkitfx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HealthConfig struct {
	Enabled bool   `default:"true"`
	Host    string `default:"0.0.0.0"`
	Port    string `default:"10010"`
}

func RunHealthServer(lc fx.Lifecycle, s fx.Shutdowner, config *Config, logger *zap.Logger, sn toolkitfx.ServiceName) {
	addr := net.JoinHostPort(config.Health.Host, config.Health.Port)
	var healthServer *grpcutils.HealthHandler
	if config.Health.Enabled {
		healthServer = grpcutils.NewHealthServer(string(sn), net.JoinHostPort(config.GRPC.Host, config.GRPC.Port), logger, addr)

		localCtx, cancel := context.WithCancel(context.Background())

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					logger.Info("starting  gRPC health server...", zap.String("addr", addr))
					err := healthServer.GServer.Serve()
					if err != nil {
						logger.Error("HealthServer gRPC server stopping down due to error", zap.Error(err))
						_ = s.Shutdown()
					}
				}()

				go func() {
					err := healthServer.CheckConnection(localCtx)
					if err != nil {
						logger.Error("HealthServer CheckConnection stopping down due to error", zap.Error(err))
						_ = s.Shutdown()
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				cancel()
				if err := healthServer.GServer.CleanUp(); err != nil {
					logger.Warn("Could not clean up Health server resources", zap.Error(err))
				}
				return nil
			},
		})
	}
}
