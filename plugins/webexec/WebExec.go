package webexec

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/webs"
	"reflect"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/plugins"
	"github.com/ecdiy/goserver/plugins/http"
	"strings"
)
func (app *Module) WebExec(ele *utils.Element) {
	we := &WebExec{Ele: ele}
	we.WebExec = reflect.ValueOf(we)
	method, mb := ele.AttrValue("Method")
	if mb && strings.ToLower(method) == "post" {
		plugins.GetGin(ele).POST(ele.MustAttr("Url"), we.run)
	} else {
		plugins.GetGin(ele).GET(ele.MustAttr("Url"), we.run)
	}
}


type WebExec struct {
	Ele     *utils.Element
	WebExec reflect.Value
}

func (we *WebExec) run(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("WebExec fail;", err)
		}
	}()
	wb := webs.NewParam(ctx)
	we.exec(wb)
	ctx.JSON(200, wb.Out)
}

func (we *WebExec) Job() {
	seelog.Info("Run Job ... ")
	we.exec(&webs.Param{Out: make(map[string]interface{}), Param: make(map[string]interface{})})
}

func (we *WebExec) exec(wb *webs.Param) {
	ns := we.Ele.AllNodes()
	for _, n := range ns {
		inputs := make([]reflect.Value, 2)
		inputs[0] = reflect.ValueOf(n)
		inputs[1] = reflect.ValueOf(wb)
		m := we.WebExec.MethodByName(n.Name())
		if m.IsValid() {
			v := m.Call(inputs)
			if len(v) == 1 && v[0].IsNil() {
				continue
			} else {
				wb.Out["Code"] = 1
				if len(v) >= 1 {
					wb.Out["Msg"] = v[0].Interface().(error).Error()
				}
				seelog.Error("error on.", n.Name())
				return
			}
		} else {
			seelog.Error("没有实现的方法:" + n.Name())
			return
		}
	}
}

//----
func (we *WebExec) Http(ele *utils.Element, wb *webs.Param) error {
	hc := &http.HCore{}
	return hc.DoHttp(ele, wb)
}

func (we *WebExec) ExecSp(ele *utils.Element, wb *webs.Param) error {
	sp := &webs.WebSp{Gpa: plugins.GetGpa(ele)}
	sp.Init()
	spName := ele.MustAttr("SpName")
	code := sp.SpExec(spName, wb)
	if code != 200 {
		seelog.Error("", spName)
	}
	return nil
}

func (we *WebExec) Param(ele *utils.Element, wb *webs.Param) error {
	for _, attr := range ele.Attrs {
		wb.Param[attr.Name()] = attr.Value
	}
	return nil
}

func (we *WebExec) Sql(ele *utils.Element, wb *webs.Param) {
	dao :=plugins. GetGpa(ele)
	dao.Exec(ele.Value)
}
