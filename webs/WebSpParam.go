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
	"github.com/ecdiy/goserver/utils"
	"strings"
	"fmt"
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
func (ws *WebSp) ParamWx(ele *utils.Element, data map[string]interface{}) { //TODO
	prefix := ele.MustAttr("Prefix")
	sqlEle := ele.Node("Sql")
	saveEle := ele.Node("SaveSql")
	saveSql := saveEle.Value
	saveParams := strings.Split(saveEle.MustAttr("Param"), ",")
	if sqlEle != nil {
		sql := sqlEle.Value
		params := strings.Split(sqlEle.MustAttr("Param"), ",")
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
				v, vb := wx[p.ParamName]
				if vb && len(v) > 0 {
					return v, 200
				}
				return ws.wxDo(saveParams, saveSql, prefix, wx["MpAppId"], wx["MpSecret"], wb, p)
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
			return ws.wxDo(saveParams, saveSql, prefix, MpAppId, MpSecret, wb, p)
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
func (ws *WebSp) wxDo(params []string, saveOpenIdSql, prefix, MpAppId, MpSecret string, wb *Param, p *SpParam) (interface{}, int) {
	jsCode := wb.String("js_code")
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + MpAppId +
		"&secret=" + MpSecret + "&js_code=" + jsCode + "&grant_type=authorization_code"
	h := &utils.Http{}
	m, e := h.GetMap(url)
	if e == nil {
		wb.Context.Set(prefix, m)
		errCode := fmt.Sprintln(m["errcode"])
		if errCode == "0" {
			var param []interface{}
			for _, p := range params {
				param = append(param, wb.String(p))
			}
			param = append(param, m["openid"])
			param = append(param, m["session_key"])
			ws.Gpa.Exec(saveOpenIdSql, param)
			v, vb := m[p.ParamName]
			if vb {
				return v, 200
			}
		}
	}
	return 401, 0
}

func (ws *WebSp) ginWk(ele *utils.Element, data map[string]interface{}, unFindCode int) {
	VerifyId, vb := ele.AttrValue("VerifyRef")
	if !vb {
		VerifyId = "Verify"
	}
	Verify := data[VerifyId].(BaseFun)
	prefix := ele.MustAttr("Prefix")
	ws.SpParamDoMap[prefix] = func(wb *Param, p *SpParam) (interface{}, int) {
		v2, b2 := wb.Context.Get(p.ParamName)
		if b2 {
			return v2, 200
		}
		Verify(wb)
		v2, b2 = wb.Context.Get(p.ParamName)
		if b2 {
			return v2, 200
		}
		if unFindCode == 401 {
			return 0, 401
		} else {
			return p.DefaultVal, 200
		}
	}
}

func (ws *WebSp) ParamGin(ele *utils.Element, data map[string]interface{}) {
	ws.ginWk(ele, data, 401)
}

func (ws *WebSp) ParamWk(ele *utils.Element, data map[string]interface{}) {
	ws.ginWk(ele, data, 200)
}
