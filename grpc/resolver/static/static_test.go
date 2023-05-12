package static_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/infraboard/mcube/grpc/resolver/static"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/grpclog"
)

var (
	SERVICE_NAME = "service_a"
	ctx          = context.Background()
)

// 提前启动服务端: go run grpc/examples/server/main.go
func TestResolver(t *testing.T) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr))

	// 连接到服务
	conn, err := grpc.DialContext(
		ctx,
		// Dial to "static://service_a"
		fmt.Sprintf("%s://%s", static.Scheme, SERVICE_NAME),
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

	fmt.Println("--- calling helloworld.Greeter/SayHello with weighted_round_robin ---")
	makeRPCs(conn, 20)
}

func callUnaryEcho(c ecpb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := ecpb.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/load_balancing")
	}
}

func init() {
	zap.DevelopmentSetup()
	store := static.GetStore()
	store.Add(SERVICE_NAME,
		static.NewTarget("127.0.0.1:50051").SetWeight(1),
		static.NewTarget("127.0.0.1:50052").SetWeight(4),
	)
}
