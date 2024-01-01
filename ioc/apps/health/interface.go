package health

import (
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	AppName = "health"
)

const (
	DEFAUL_HEALTH_PATH = "/healthz"
)

type HealthCheck struct {
	Path string `json:"path" yaml:"path" toml:"path" env:"HTTP_HEALTH_CHECK_PATH"`
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func NewHealthCheckRequest() *healthgrpc.HealthCheckRequest {
	return &healthgrpc.HealthCheckRequest{}
}

func NewHealth(hc *healthgrpc.HealthCheckResponse) *HealthCheckResponse {
	return &HealthCheckResponse{
		Status: hc.Status.String(),
	}
}
