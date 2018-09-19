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
	sql string
	g   *gpa.Gpa
}

func GetRpc(dao *gpa.Gpa, sql string) *RpcUser {
	return &RpcUser{g: dao, sql: sql}
}

func (s *RpcUser) Verify(c context.Context, in *Token) (*UserBase, error) {
	ub := &UserBase{}
	if len(in.Token) > 1 {
		v, b := TokenMap[in.Token]
		if b {
			ub.UserId = v
			setUb(ub)
		} else {
			m, b, ee := s.g.QueryMapStringString(s.sql, in.Token, in.Ua)
			if ee == nil && b {
				uId, _ := strconv.ParseInt(m["UserId"], 10, 0)
				TokenMap[in.Token] = uId
				setUb(ub)
			} else {
				ub.Result = false
			}
		}
	}
	return ub, nil
}

func setUb(ub *UserBase) {
	m, mb := UserMap[ub.UserId]
	if mb {
		ub.Result = true
		uId, _ := strconv.ParseInt(m["UserId"], 10, 0)
		un, unb := m["Username"]
		if unb {
			ub.Username = un
		}
		if len(m) > 2 {
			be, _ := json.Marshal(m)
			ub.AppendJson = string(be)
		}
		ub.UserId = uId
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
func RegisterSpAjax(g *gpa.Gpa, eng *gin.Engine, rpc *RpcUser, RpcUserHost, url, ext, tokenName string) {
	auth := func(gCtx *gin.Context) (bool, int64) {
		gCtx.Set("CallAuth", 1)
		wb := WebBaseNew(gCtx)
		tokenVal := wb.String(tokenName)
		if len(tokenVal) > 1 {
			auth := false
			var ub *UserBase
			if len(RpcUserHost) > 1 {
				rpcUser(RpcUserHost, func(client RpcUserClient, ctx context.Context) {
					ub, _ = client.Verify(ctx, &Token{Token: tokenVal, Ua: wb.Ua})
				})
			} else {
				ub, _ = rpc.Verify(nil, &Token{Token: tokenVal, Ua: wb.Ua})

			}
			if ub.Result {
				gCtx.Set("UserId", ub.UserId)
				gCtx.Set("Username", ub.Username)
				if len(ub.AppendJson) > 1 {
					var data map[string]interface{}
					je := json.Unmarshal([]byte(ub.AppendJson), &data)
					if je == nil {
						for k, v := range data {
							gCtx.Set(k, v)
						}
					}
				}
				auth = true
			} else {
				auth = false
			}
			return auth, ub.UserId
		} else {
			return false, 0
		}
	}
	if !gin.IsDebugging() {
		spInitCache(g, auth, ext)
	}
	eng.POST(url, func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error("sp un catch error;", err)
			}
		}()
		sp(g, c, ext, auth)
	})
}
