package sp

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/webs"
	"github.com/ecdiy/goserver/plugins"
)

func init() {

}
func Sp(ele *utils.Element) {
	sp := &webs.WebSp{Gpa: plugins.GetGpa(ele), Engine: plugins.GetGin(ele)}
	sp.Init()
	doSubElement(ele, sp)
	put(ele, sp)
}
