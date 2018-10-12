package webs

import (
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
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

// "/sp/:sp"   "Ajax"
// "/spa/:sp"  "Admin"
func RegisterSpAjax(auth func(param *Param) *UserBase, url, ext string) {
	if !gin.IsDebugging() {
		spInitCache(Gpa, auth, ext)
	}
	Gin.POST(url, func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error("sp un catch error;", url, ";", err)
			}
		}()
		spName := c.Param("sp") + ext
		wb := NewParam(c)
		code := SpExec(spName, Gpa, wb, auth)
		if code == 200 {
			c.JSON(200, wb.Out)
		} else {
			seelog.Error("数据存储过程错误:"+spName, ";", code)
			c.AbortWithStatus(code)
		}
	})
}

func RegisterReload(url string) {
	Gin.GET(url, func(i *gin.Context) {
		spCache = make(map[string]*Sp)
		i.String(200, "clear cache ok.")
	})
}
