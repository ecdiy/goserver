package webs

import (
	"encoding/json"
	"context"
	"strconv"
	"github.com/ecdiy/goserver/gpa"
	"google.golang.org/grpc"
	"net"
	"github.com/cihub/seelog"
	"fmt"
	"github.com/ecdiy/goserver/utils"
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

func NewVerify(ele *utils.Element, Gpa *gpa.Gpa, putFunRun func(fun func())) BaseFun {
	v := &Verify{
		tkName: ele.MustAttr("TokenName"),
	}
	v.ResultFlagName, v.ResultFlagNameExt = ele.AttrValue("ResultFlagName")
	v.sql, v.sqlExt = ele.AttrValue("Sql")
	if !v.sqlExt {
		v.sql = ele.Value
		if len(strings.TrimSpace(v.sql)) > 10 {
			v.sqlExt = true
		}
	}
	v.RpcHost, v.rhExt = ele.AttrValue("RpcHost")
	if v.sqlExt {
		v.ver = &RpcUser{Sql: v.sql, Gpa: Gpa}
		if v.rhExt {
			putFunRun(func() {
				RpcRegister(v.RpcHost, func(s *grpc.Server) {
					RegisterRpcUserServer(s, v.ver)
				})
			})
		} else {
			seelog.Warn("没有设置RpcHost,无须启动绑定Rpc")
		}
	} else {
		if !v.rhExt {
			panic("参数设置错误[Sql,Sql+RpcHost,RpcHost]三种方式")
		}
	}
	return v.DoAuth
}

type Verify struct {
	tkName, RpcHost, ResultFlagName, sql string
	sqlExt, rhExt, ResultFlagNameExt     bool
	ver                                  *RpcUser
}

func (v *Verify) DoAuth(wb *Param, ps ... interface{}) interface{} {
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
	tk := &Token{Token: wb.String(v.tkName), Ua: wb.Ua}
	if v.sqlExt {
		ub, _ = v.ver.Verify(nil, tk)
	} else {
		rpcUser(v.RpcHost, func(client RpcUserClient, ctx context.Context) {
			ub, _ = client.Verify(ctx, tk)
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
	} else {
		if utils.EnvIsDev {
			seelog.Warn("认证失败，token-name=", v.tkName, ":", tk.Token, ";", tk.Ua)
		}
	}
	if v.ResultFlagNameExt {
		wb.Context.Set(v.ResultFlagName, ub != nil && ub.Result)
	}
	return ub
}
