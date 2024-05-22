package fxgrpc

import (
	"context"
	"net"
	"net/http"
	"time"

	grpcutils "github.com/ecumenos-social/grpc-utils"
	"github.com/ecumenos-social/toolkitfx"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	GRPCHost string `default:"0.0.0.0"`
	GRPCPort string `default:"8080"`

	HTTPGatewayHost string `default:"0.0.0.0"`
	HTTPGatewayPort string `default:"9090"`

	EnabledHealthServer bool   `default:"true"`
	HealthServerHost    string `default:"0.0.0.0"`
	HealthServerPort    string `default:"10010"`

	LivenessGatewayHost string `default:"0.0.0.0"`
	LivenessGatewayPort string `default:"8086"`
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

func RunHealthServer(lc fx.Lifecycle, s fx.Shutdowner, cfg *Config, logger *zap.Logger, sn toolkitfx.ServiceName) {
	addr := net.JoinHostPort(cfg.HealthServerHost, cfg.HealthServerPort)
	var healthServer *grpcutils.HealthHandler
	if cfg.EnabledHealthServer {
		healthServer = grpcutils.NewHealthServer(string(sn), net.JoinHostPort(cfg.GRPCHost, cfg.GRPCPort), logger, addr)

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

type GatewayHandler struct {
	Handler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

func RunHTTPGateway(lc fx.Lifecycle, s fx.Shutdowner, logger *zap.Logger, cfg *Config, g *GatewayHandler) error {
	addr := net.JoinHostPort(cfg.HTTPGatewayHost, cfg.HTTPGatewayPort)
	mux := runtime.NewServeMux()
	conn := grpcutils.NewClientConnection(net.JoinHostPort(cfg.GRPCHost, cfg.GRPCPort))

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

type LivenessHandler struct {
	Handler http.Handler
}

func RunLivenessGateway(lc fx.Lifecycle, s fx.Shutdowner, logger *zap.Logger, cfg *Config, h *LivenessHandler) error {
	addr := net.JoinHostPort(cfg.LivenessGatewayHost, cfg.LivenessGatewayPort)
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
