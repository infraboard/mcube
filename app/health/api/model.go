package api

import (
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func NewHealth(hc *healthgrpc.HealthCheckResponse) *Health {
	return &Health{
		Status: hc.Status.String(),
	}
}

type Health struct {
	Status string `json:"status"`
}
