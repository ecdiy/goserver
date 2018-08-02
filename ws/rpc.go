package ws

import (
	"google.golang.org/grpc"

	"net"
	"github.com/cihub/seelog"
	"context"
)

func RpcUser(fun func(client RpcUserClient, ctx context.Context)) {
	conn, err := grpc.DialContext(context.Background(), RpcUserHost, grpc.WithInsecure())
	if err == nil {
		defer conn.Close()
		client := NewRpcUserClient(conn)
		fun(client, context.Background())
	}
}

func RpcRegister(addr string, regFunc func(server *grpc.Server)) {
	s := grpc.NewServer()
	rh := EnvParam("RpcHost")
	if rh == "" {
		rh = addr
	}
	lis, err := net.Listen("tcp", rh)
	if err == nil {
		regFunc(s)
		seelog.Info("grpc bind: " + rh)
		s.Serve(lis)
	} else {
		panic("启动微服务失败")
	}
}
