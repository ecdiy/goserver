package webs

import (
	"encoding/json"
	"context"
	"strconv"
	"utils/gpa"
	"google.golang.org/grpc"
	"utils"
	"net"
	"github.com/cihub/seelog"
	"fmt"
	"utils/xml"
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

func ubx(ub *UserBase, param *Param) {
	if ub.Result {
		param.Context.Set("UserId", ub.UserId)
		param.Context.Set("Username", ub.Username)
		if len(ub.AppendJson) > 1 {
			var data map[string]interface{}
			je := json.Unmarshal([]byte(ub.AppendJson), &data)
			if je == nil {
				for k, v := range data {
					param.Context.Set(k, v)
				}
			}
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

func NewBaseFun(ele *xml.Element, Gpa *gpa.Gpa) BaseFun {
	RpcHost, rhb := ele.AttrValue("RpcHost")
	tkName := ele.MustAttr("TokenName")
	var ver *RpcUser
	if !rhb {
		ver = &RpcUser{Sql: ele.MustAttr("Sql"), Gpa: Gpa}
	}
	return func(wb *Param, ps ... interface{}) interface{} {
		_, ext := wb.Context.Get(VerifyCallFlag)
		if ext {
			v, vb := wb.Context.Get("Verify")
			if vb {
				return v.(*UserBase)
			}
			seelog.Error("gin context not save UserBase")
			return nil
		}
		wb.Context.Set(VerifyCallFlag, true)
		var ub *UserBase
		if rhb {
			rpcUser(RpcHost, func(client RpcUserClient, ctx context.Context) {
				ub, _ = client.Verify(ctx, &Token{Token: wb.String(tkName), Ua: wb.Ua})
			})
		} else {
			ub, _ = ver.Verify(nil, &Token{Token: wb.String(tkName), Ua: wb.Ua})
		}
		ubx(ub, wb)
		return ub
	}
}
