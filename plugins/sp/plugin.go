package sp

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins/verify"
	"strings"
	"fmt"
	"github.com/ecdiy/goserver/plugins"
)

const (
	SqlSpAll  = "select name,CONVERT(param_list USING utf8) param_list,`comment` from mysql.proc c where db=DATABASE() and `type`='PROCEDURE'"
	SqlSpInfo = "select name,CONVERT(param_list USING utf8) param_list,`comment` from mysql.proc c where db=DATABASE() and `type`='PROCEDURE' and name=?"
)

var (
	spReloadFun = make(map[string]func(c *utils.Param) *verify.UserBase)

	ParamDoMap  = make(map[string]ParamValFunc)
	//存储过程参数处理规制
)

func FmtSql(sql string, param ... interface{}) string {
	s := sql
	if param != nil {
		for _, p := range param {
			s = strings.Replace(s, "?", "'"+fmt.Sprint(p)+"'", 1)
		}
	}
	return s
}

func init() {
	plugins.RegisterPlugin("SpParamGin", func(ele *utils.Element) interface{} {
		vf := plugins.GetRef(ele, "Verify")
		ParamDoMap[ ele.Attr("Prefix", "gin")] = ginWk(vf.(plugins.BaseFun), 401)
		return nil
	})
	plugins.RegisterPlugin("SpParamWk", func(ele *utils.Element) interface{} {
		vf := plugins.GetRef(ele, "Verify")
		ParamDoMap[ ele.MustAttr("Prefix")] = ginWk(vf.(plugins.BaseFun), 200)
		return nil
	})
	plugins.RegisterPlugin("SpParamWx", func(ele *utils.Element) interface{} {
		ParamDoMap[ ele.Attr("Prefix", "wx")] = ParamWx(ele)
		return nil
	})
}
