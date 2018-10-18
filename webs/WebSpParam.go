package webs

/**
存储过程参数规则
1.ua
2.in*
3.wx*
4.gin*
5.wk*
 */

import (
	"github.com/cihub/seelog"
	"utils/xml"
	"context"
	"utils/http"
	"strings"
)

//--
func ParamUa(wb *Param, p *SpParam) (interface{}, int) {
	return wb.Ua, 200
}

func ParamIn(ctx *Param, p *SpParam) (interface{}, int) {
	v := ctx.String(p.ParamName)
	if v == "" {
		v = p.DefaultVal
	}
	return v, 200
}

/*--微信获取用户Id
1. db appId,secret.
2. http get openId
3. query openId to userId.
*/
func (ws *SpWeb) ParamWx(ele *xml.Element, data map[string]interface{}) {
	prefix := ele.MustAttr("Prefix")
	sql, sqlExt := ele.AttrValue("Sql")
	if sqlExt {
		SqlParam := ele.MustAttr("SqlParam")
		params := strings.Split(SqlParam, ",")
		ws.SpParamDoMap[prefix] = func(wb *Param, p *SpParam) (interface{}, int) {
			v, vc := wxv(prefix, wb, p)
			if vc != 0 {
				return v, vc
			}
			var param []interface{}
			for _, p := range params {
				param = append(param, wb.String(p))
			}
			wx, fd, _ := ws.Gpa.QueryMapStringString(sql, param...)
			if fd {
				return wxDo(prefix, wx["MpAppId"], wx["MpSecret"], wb, p)
			} else {
				seelog.Error("没有匹配上AppId,Secret,param=", param)
			}
			return 401, 0
		}
	} else {
		MpAppId := ele.MustAttr("MpAppId")
		MpSecret := ele.MustAttr("MpSecret")
		ws.SpParamDoMap[prefix] = func(wb *Param, p *SpParam) (interface{}, int) {
			v, vc := wxv(prefix, wb, p)
			if vc != 0 {
				return v, vc
			}
			return wxDo(prefix, MpAppId, MpSecret, wb, p)
		}
	}
}

func wxv(prefix string, wb *Param, p *SpParam) (interface{}, int) {
	v, vb := wb.Context.Get(prefix)
	if vb {
		pv, pvb := v.(map[string]interface{})[p.ParamName]
		if pvb {
			return pv, 200
		}
		return "", 401
	}
	return "", 0
}

// https://api.weixin.qq.com/sns/jscode2session?appid=wx1d88c080ac473555&secret=0d19956a54b1ce949902f5f7d7523658&js_code=0011AwBI1DMSX600BeEI1FDNBI11AwB6&grant_type=authorization_code
func wxDo(prefix, MpAppId, MpSecret string, wb *Param, p *SpParam) (interface{}, int) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + MpAppId +
		"&secret=" + MpSecret +
		"&js_code=" + wb.String("js_code") + "&grant_type=authorization_code"
	h := &http.Http{}
	m, e := h.GetMap(url)
	if e == nil {
		wb.Context.Set(prefix, m)
	}
	v, vb := m[p.ParamName]
	if vb {
		return v, 200
	}
	return 401, 0
}

func (ws *SpWeb) ginWk(ele *xml.Element, data map[string]interface{}, unFindCode int) {
	tkName := ele.MustAttr("TokenName")
	prefix := ele.MustAttr("Prefix")
	RpcHost, rhb := ele.AttrValue("RpcHost")
	if rhb {
		ws.SpParamDoMap[prefix] = func(wb *Param, p *SpParam) (interface{}, int) {
			v, b := wb.Context.Get(p.ParamName)
			if b {
				return v, 200
			}
			var ub *UserBase
			rpcUser(RpcHost, func(client RpcUserClient, ctx context.Context) {
				ub, _ = client.Verify(ctx, &Token{Token: wb.String(tkName), Ua: wb.Ua})
			})
			if ub.Result {
				ubx(ub, wb)
				v2, b2 := wb.Context.Get(p.ParamName)
				if b2 {
					return v2, 200
				}
			}
			seelog.Error("ctx.Get not find.", p.ParamName)
			if unFindCode == 401 {
				return 0, 401
			} else {
				return p.DefaultVal, 200
			}
		}
		return
	}
	RpcRef, rr := ele.AttrValue("RpcRef")
	var ver *RpcUser
	if rr {
		ver = data[RpcRef].(*RpcUser)
	} else {
		ver = &RpcUser{Sql: ele.MustAttr("Sql"), Gpa: ws.Gpa}
	}
	ws.SpParamDoMap[prefix] = func(wb *Param, p *SpParam) (interface{}, int) {
		v, b := wb.Context.Get(p.ParamName)
		if b {
			return v, 200
		}
		ver.Verify(nil, &Token{Token: wb.String(tkName), Ua: wb.Ua})
		v2, b2 := wb.Context.Get(p.ParamName)
		if b2 {
			return v2, 200
		}
		seelog.Error("ctx.Get not find.", p.ParamName)
		if unFindCode == 401 {
			return 0, 401
		} else {
			return p.DefaultVal, 200
		}
	}
}

func (ws *SpWeb) ParamGin(ele *xml.Element, data map[string]interface{}) {
	ws.ginWk(ele, data, 401)
}

func (ws *SpWeb) ParamWk(ele *xml.Element, data map[string]interface{}) {
	ws.ginWk(ele, data, 200)
}
