package static_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcube/grpc/resolver/static"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	SERVICE_NAME = "service_a"
	ctx          = context.Background()
)

func TestResolver(t *testing.T) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// 连接到服务
	conn, err := grpc.DialContext(
		ctx,
		// Dial to "static://service_a"
		fmt.Sprintf("%s://%s", static.Scheme, SERVICE_NAME),
		// 认证
		grpc.WithPerRPCCredentials(rpc.NewAuthenticationFromEnv()),
		// 不开启TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// gprc 支持的负载均衡策略: https://github.com/grpc/grpc/blob/master/doc/load-balancing.md
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"weighted_round_robin":{}}]}`),
		// 直到建立连接
		grpc.WithBlock(),
	)

	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
}

func init() {
	zap.DevelopmentSetup()
	store := static.GetStore()
	store.Add(SERVICE_NAME,
		static.NewTarget("10.10.10.10"),
		static.NewTarget("10.10.10.11"),
		static.NewTarget("10.10.10.12"),
	)
}
