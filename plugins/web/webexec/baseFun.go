package webexec

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins/http"
	"github.com/ecdiy/goserver/plugins/sp"
	"github.com/ecdiy/goserver/plugins"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/gpa"
)

//----
func (we *WebExec) Http(ele *utils.Element, wb *utils.Param) error {
	hc := &http.HCore{}
	return hc.DoHttp(ele, wb)
}

func (we *WebExec) Sp(ele *utils.Element, wb *utils.Param) error {
	sp := &sp.WebSp{Gpa: plugins.GetRef(ele, "Gpa").(*gpa.Gpa)}
	sp.Init(ele)
	spName := ele.MustAttr("SpName")
	code := sp.SpExec(spName, wb)
	if code != 200 {
		seelog.Error("", spName)
	}
	return nil
}

func (we *WebExec) Param(ele *utils.Element, wb *utils.Param) error {
	for _, attr := range ele.Attrs {
		wb.Param[attr.Name()] = attr.Value
	}
	return nil
}

func (we *WebExec) Sql(ele *utils.Element, wb *utils.Param) {
	dao := plugins.GetRef(ele, "Gpa").(*gpa.Gpa)
	dao.Exec(ele.Value)
}
