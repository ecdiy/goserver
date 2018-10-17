package webs

import (
	"google.golang.org/grpc"
	"github.com/cihub/seelog"
	"context"
	"utils"
	"net"
	"fmt"
)

var (
	TokenMap = make(map[string]int64)
	UserMap  = make(map[int64]map[string]string)
)

func RpcRegister(addr string, regFunc ... func(server *grpc.Server)) {
	s := grpc.NewServer()
	rh := utils.EnvParam("RpcHost")
	if rh == "" {
		rh = addr
	}
	lis, err := net.Listen("tcp", rh)
	if err == nil {
		for _, v := range regFunc {
			v(s)
		}
		//regFunc(s)
		seelog.Info("grpc bind: " + rh)
		s.Serve(lis)
	} else {
		panic("启动微服务失败:" + fmt.Sprintln(err) + ";Input Addr=" + addr + ";RpcHost=" + rh)
	}
}

func rpcUser(RpcUserHost string, fun func(client RpcUserClient, ctx context.Context)) {
	conn, err := grpc.DialContext(context.Background(), RpcUserHost, grpc.WithInsecure())
	if err == nil {
		defer conn.Close()
		client := NewRpcUserClient(conn)
		fun(client, context.Background())
	}
}
