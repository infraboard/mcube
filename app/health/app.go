package health

import (
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	AppName = "health"
)

func NewHealthCheckRequest() *healthgrpc.HealthCheckRequest {
	return &healthgrpc.HealthCheckRequest{}
}
