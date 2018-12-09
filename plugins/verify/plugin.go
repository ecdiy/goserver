package verify

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/gpa"
)

func init() {
	plugins.RegisterPlugin("Verify", func(ele *utils.Element) interface{} {
		vb := NewVerify(ele, plugins.GetRef(ele, "Gpa").(*gpa.Gpa), plugins.PutFunRun)
		//put(ele, vb)
		tfn, ext := ele.AttrValue("TplFunName")
		if ext {
			//webs.RegisterBaseFun(tfn, vb)
			seelog.Info("添加Verify函数:", tfn)
		}
		return vb
	})
}
