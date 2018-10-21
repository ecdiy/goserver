package webs

import (
	"encoding/json"
	"context"
	"strconv"
	"goserver/gpa"
	"google.golang.org/grpc"
	"goserver"
	"net"
	"github.com/cihub/seelog"
	"fmt"
	"goserver/xml"
	"strings"
)

type RpcUser struct {
	Sql string
	Gpa *gpa.Gpa
}

func (s *RpcUser) Verify(c context.Context, in *Token) (*UserBase, error) {
	ub := &UserBase{}
	if len(in.Token) > 1 {
		m, b, ee := s.Gpa.QueryMapStringString(s.Sql, in.Token, in.Ua)
		if ee == nil && b {
			uId, _ := strconv.ParseInt(m["UserId"], 10, 0)
			ub.UserId = uId
			ub.Result = true
			if len(m) > 1 {
				be, _ := json.Marshal(m)
				ub.AppendJson = string(be)
			}
		} else {
			ub.Result = false
		}
	}
	return ub, nil
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

func NewVerify(ele *xml.Element, Gpa *gpa.Gpa, putFunRun func(fun func())) BaseFun {
	sql, sqlExt := ele.AttrValue("Sql")
	if !sqlExt {
		sql = ele.Value
		if len(strings.TrimSpace(sql)) > 10 {
			sqlExt = true
		}
	}
	RpcHost, rhExt := ele.AttrValue("RpcHost")
	tkName := ele.MustAttr("TokenName")
	ResultFlagName, ResultFlagNameExt := ele.AttrValue("ResultFlagName")
	var ver *RpcUser
	if sqlExt {
		ver = &RpcUser{Sql: sql, Gpa: Gpa}
		if rhExt {
			putFunRun(func() {
				RpcRegister(RpcHost, func(s *grpc.Server) {
					RegisterRpcUserServer(s, ver)
				})
			})
		} else {
			seelog.Warn("没有设置RpcHost,无须启动绑定Rpc")
		}
	} else {
		if !rhExt {
			panic("参数设置错误[Sql,Sql+RpcHost,RpcHost]三种方式")
		}
	}
	return func(wb *Param, ps ... interface{}) interface{} {
		_, ext := wb.Context.Get(VerifyCallFlag)
		if ext {
			v, vb := wb.Context.Get("Verify")
			if vb {
				return v.(*UserBase)
			}
			seelog.Error("gin context not save UserBase.Verify fail?")
			return nil
		}
		wb.Context.Set(VerifyCallFlag, true)
		var ub *UserBase
		if sqlExt {
			ub, _ = ver.Verify(nil, &Token{Token: wb.String(tkName), Ua: wb.Ua})
		} else {
			rpcUser(RpcHost, func(client RpcUserClient, ctx context.Context) {
				ub, _ = client.Verify(ctx, &Token{Token: wb.String(tkName), Ua: wb.Ua})
			})
		}
		if ub != nil && ub.Result {
			if len(ub.AppendJson) > 1 {
				var data map[string]interface{}
				je := json.Unmarshal([]byte(ub.AppendJson), &data)
				if je == nil {
					for k, v := range data {
						wb.Context.Set(k, v)
					}
				}
			} else {
				wb.Context.Set("UserId", ub.UserId)
			}
		}
		if ResultFlagNameExt {
			wb.Context.Set(ResultFlagName, ub != nil && ub.Result)
		}
		return ub
	}
}
