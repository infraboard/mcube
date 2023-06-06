package impl

import (
	"context"

	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func (i *impl) Check(ctx context.Context, req *healthgrpc.HealthCheckRequest) (
	*healthgrpc.HealthCheckResponse, error) {
	return i.Server.Check(ctx, req)
}
