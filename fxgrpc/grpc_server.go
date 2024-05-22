package fxgrpc

import (
	"context"
	"time"

	grpcutils "github.com/ecumenos-social/grpc-utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type GRPCConfig struct {
	Host                                    string        `default:"0.0.0.0"`
	Port                                    string        `default:"8080"`
	MaxConnectionAge                        time.Duration `default:"5m"`
	KeepAliveEnforcementMinTime             time.Duration `default:"1m"`
	KeepAliveEnforcementPermitWithoutStream bool          `default:"true"`
}

type GRPCServer struct {
	Server *grpcutils.Server
}

func RunRegisteredGRPCServer(lc fx.Lifecycle, s fx.Shutdowner, logger *zap.Logger, gs *GRPCServer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("starting gRPC server...")
				if err := gs.Server.Serve(); err != nil {
					logger.Error("registered gRPC server stopping down due to error", zap.Error(err))
					_ = s.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := gs.Server.CleanUp(); err != nil {
				logger.Error("can not clean up server resources", zap.Error(err))
			}
			return nil
		},
	})
}
