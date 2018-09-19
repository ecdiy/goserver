package webs

import (
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"context"
	"utils/gpa"
	"encoding/json"
	"strconv"
)

var (
	TokenMap = make(map[string]map[string]string)
)

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
func Register(g *gpa.Gpa, eng *gin.Engine, RpcUserHost, url, ext, tokenName, sql string) {
	auth := func(gCtx *gin.Context) (bool, int64) {
		gCtx.Set("CallAuth", 1)
		wb := WebBaseNew(gCtx)
		tokenVal := wb.String(tokenName)
		if len(tokenVal) > 1 {
			userId := int64(0)
			auth := false
			if len(RpcUserHost) > 1 {
				rpcUser(RpcUserHost, func(client RpcUserClient, ctx context.Context) {
					sc, _ := client.Verify(ctx, &Token{Token: tokenVal, Ua: wb.Ua})
					if sc.Result {
						userId = sc.UserId
						gCtx.Set("UserId", userId)
						gCtx.Set("Username", sc.Username)
						if len(sc.AppendJson) > 1 {
							var data map[string]interface{}
							je := json.Unmarshal([]byte(sc.AppendJson), &data)
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
				})
				return auth, userId
			} else {
				uKey := ext + "." + tokenName
				v, b := TokenMap[uKey]
				if b {
					for k, vv := range v {
						gCtx.Set(k, vv)
					}
					uId, _ := strconv.ParseInt(v["UserId"], 10, 0)
					return true, uId
				} else {
					m, b, ee := g.QueryMapStringString(sql, tokenVal, wb.Ua)
					if ee == nil && b {
						TokenMap[uKey] = m
						uId, _ := strconv.ParseInt(v["UserId"], 10, 0)
						return true, uId
					}
					return false, 0
				}
			}
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
