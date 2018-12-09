package sp

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"strings"
	"github.com/ecdiy/goserver/plugins/sp"
)

func HandleBase(ele *utils.Element) func(c *gin.Context) {
	ws := &sp.WebSp{}
	ws.Init(ele)
	sh := &Handle{ws: ws, spSuffix: ele.MustAttr("SpSuffix")}

	if !utils.EnvIsDev {
		list, err := ws.Gpa.ListArrayString(sp.SqlSpAll)
		if err != nil {
			seelog.Error("查询所有的存储过程出错：", err)
		} else {
			if list != nil {
				sps := ""
				for _, val := range list {
					if strings.LastIndex(val[0], sh.spSuffix) == len(val[0])-len(sh.spSuffix) {
						sp, b := ws.NewSp(val)
						if b {
							ws.SpCache [sp.Name] = sp
							sps += sp.Sql + ","
						}
					}
				}
				seelog.Info("~~~~;\n\t", sps)
			}
		}
	}
	sh.RuleSp, sh.rule = ele.AttrValue("RuleSp")
	return sh.Handle
	//reloadUrl, rExt := ele.AttrValue("ReloadUrl")
	//if rExt {
	//	ws.Engine.GET(reloadUrl, func(i *gin.Context) {
	//		ws.SpCache = make(map[string]*Sp)
	//		i.String(200, "clear cache ok.")
	//	})
	//}
}
