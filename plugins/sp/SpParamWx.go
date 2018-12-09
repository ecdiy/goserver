package sp

import (
	"github.com/ecdiy/goserver/utils"
	"fmt"
	"strings"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/plugins"
	"github.com/ecdiy/goserver/gpa"
)

/*--微信获取用户Id
1. db appId,secret.
2. http get openId
3. query openId to userId.
*/

func ParamWx(ele *utils.Element) func(wb *utils.Param, p *SpParam) (interface{}, int) { //TODO
	Gpa := plugins.GetRef(ele, "Gpa").(*gpa.Gpa)
	prefix := ele.MustAttr("Prefix")
	sqlEle := ele.Node("Sql")
	saveEle := ele.Node("SaveSql")
	saveSql := saveEle.Value
	saveParams := strings.Split(saveEle.MustAttr("Param"), ",")
	if sqlEle != nil {
		sql := sqlEle.Value
		params := strings.Split(sqlEle.MustAttr("Param"), ",")
		return func(wb *utils.Param, p *SpParam) (interface{}, int) {
			v, vc := wxv(prefix, wb, p)
			if vc != 0 {
				return v, vc
			}
			var param []interface{}
			for _, p := range params {
				param = append(param, wb.String(p))
			}
			wx, fd, _ := Gpa.QueryMapStringString(sql, param...)
			if fd {
				v, vb := wx[p.ParamName]
				if vb && len(v) > 0 {
					return v, 200
				}
				return wxDo(Gpa, saveParams, saveSql, prefix, wx["MpAppId"], wx["MpSecret"], wb, p)
			} else {
				seelog.Error("没有匹配上AppId,Secret,param=", param)
			}
			return 401, 0
		}
	} else {
		MpAppId := ele.MustAttr("MpAppId")
		MpSecret := ele.MustAttr("MpSecret")
		return func(wb *utils.Param, p *SpParam) (interface{}, int) {
			v, vc := wxv(prefix, wb, p)
			if vc != 0 {
				return v, vc
			}
			return wxDo(Gpa, saveParams, saveSql, prefix, MpAppId, MpSecret, wb, p)
		}
	}
}

// https://api.weixin.qq.com/sns/jscode2session?appid=wx1d88c080ac473555&secret=0d19956a54b1ce949902f5f7d7523658&js_code=0011AwBI1DMSX600BeEI1FDNBI11AwB6&grant_type=authorization_code
func wxDo(Gpa *gpa.Gpa, params []string, saveOpenIdSql, prefix, MpAppId, MpSecret string, wb *utils.Param, p *SpParam) (interface{}, int) {
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
			Gpa.Exec(saveOpenIdSql, param)
			v, vb := m[p.ParamName]
			if vb {
				return v, 200
			}
		}
	}
	return 401, 0
}

func wxv(prefix string, wb *utils.Param, p *SpParam) (interface{}, int) {
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
