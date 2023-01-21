package main

import (
	"flag"
	"fmt"

	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/internal/config"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/internal/server"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/internal/svc"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/types/notifications/v1alpha1"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/notificationapi.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		v1alpha1.RegisterNotificationAPIServiceServer(grpcServer, server.NewNotificationAPIServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
