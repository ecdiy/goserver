package webs

import (
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"context"
	"utils/gpa"
	"encoding/json"
	"strconv"
	"utils"
	"net"
)

var (
	TokenMap = make(map[string]int64)
	UserMap  = make(map[int64]map[string]string)
)

type RpcUser struct {
	Sql string
	Gpa *gpa.Gpa
}

func (s *RpcUser) Verify(c context.Context, in *Token) (*UserBase, error) {
	ub := &UserBase{}
	if len(in.Token) > 1 {
		v, b := TokenMap[in.Token]
		if b {
			ub.UserId = v
			setUb(ub)
		} else {
			m, b, ee := s.Gpa.QueryMapStringString(s.Sql, in.Token, in.Ua)
			if ee == nil && b {
				uId, _ := strconv.ParseInt(m["UserId"], 10, 0)
				TokenMap[in.Token] = uId
				ub.UserId = uId
				setUb(ub)
			} else {
				ub.Result = false
			}
		}
	}
	return ub, nil
}

func setUb(ub *UserBase) {
	ub.Result = true
	m, mb := UserMap[ub.UserId]
	if mb {
		un, unb := m["Username"]
		if unb {
			ub.Username = un
		}
		if len(m) > 2 {
			be, _ := json.Marshal(m)
			ub.AppendJson = string(be)
		}
	}
}

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
		panic("启动微服务失败")
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
func RegisterSpAjax(g *gpa.Gpa, eng *gin.Engine, auth func(param *Param) *UserBase, url, ext string) {
	if !gin.IsDebugging() {
		spInitCache(g, auth, ext)
	}
	eng.POST(url, func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error("sp un catch error;", url, ";", err)
			}
		}()
		spName := c.Param("sp") + ext
		wb := NewParam(c)
		code := SpExec(spName, g, wb, auth)
		if code == 200 {
			c.JSON(200, wb.Out)
		} else {
			seelog.Error("数据存储过程错误:"+spName, ";", code)
			c.AbortWithStatus(code)
		}
	})
}

func RegisterReload(url string, eng *gin.Engine) {
	eng.GET(url, func(i *gin.Context) {
		spCache = make(map[string]*Sp)
		i.String(200, "clear cache ok.")
	})
}
