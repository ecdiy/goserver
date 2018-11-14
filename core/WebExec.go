package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/webs"
	"reflect"
	"github.com/cihub/seelog"
)

type WebExec struct {
	ele     *utils.Element
	webExec reflect.Value
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

func (we *WebExec) job() {
	we.exec(&webs.Param{Out: make(map[string]interface{}), Param: make(map[string]interface{})})
}

func (we *WebExec) exec(wb *webs.Param) {
	ns := we.ele.AllNodes()
	for _, n := range ns {
		inputs := make([]reflect.Value, 2)
		inputs[0] = reflect.ValueOf(n)
		inputs[1] = reflect.ValueOf(wb)
		m := we.webExec.MethodByName(n.Name())
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
	hc := &HttpCore{}
	return hc.DoHttp(ele, wb)
}

func (we *WebExec) ExecSp(ele *utils.Element, wb *webs.Param) error {
	sp := &webs.WebSp{Gpa: getGpa(ele)}
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
	dao := getGpa(ele)
	dao.Exec(ele.Value)
}
