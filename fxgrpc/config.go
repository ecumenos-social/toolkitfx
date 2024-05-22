package fxgrpc

type Config struct {
	GRPC            GRPCConfig
	Health          HealthConfig
	HTTPGateway     HTTPGatewayConfig
	LivenessGateway LivenessGatewayConfig
}
