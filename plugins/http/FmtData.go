package http

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/cihub/seelog"
	"strings"
	"github.com/ecdiy/goserver/plugins"
	"regexp"
	"strconv"
	"fmt"
	"reflect"
	"github.com/ecdiy/goserver/gpa"
	"github.com/ecdiy/goserver/plugins/sp"
)

type FmtData struct {
	items    []string
	Dao      *gpa.Gpa
	ErrorReg string
}

func (fd *FmtData) Spit(ele *utils.Element, html string, param *utils.Param) {
	Spit := ele.Node("SpitString")
	if Spit == nil {
		seelog.Error("<SpitString>没有设置")
		return
	}
	fd.items = strings.Split(html, Spit.Value)

	if fd.items != nil {
		Param := ele.Node("Param")
		Sp, spExt := Param.AttrValue("Sp")

		sp := &sp.WebSp{Gpa: plugins.GetRef(ele, "Gpa").(*gpa.Gpa)}
		sp.Init(ele)

		ItemInclude, ItemIncludeExt := ele.AttrValue("ItemInclude")
		//for _, it := range fd.items {
		//errRes:=""
		saveNode := 0
		errorNode := 0
		for i := len(fd.items) - 1; i >= 0; i-- {
			it := fd.items[i]
			if ItemIncludeExt {
				if strings.Index(it, ItemInclude) < 0 {
					continue
				}
			}
			val := fd.getParam(it, Param)
			if spExt && len(val) > 0 {
				wb := &utils.Param{Out: make(map[string]interface{}), Param: val}
				for k, v := range param.Param {
					wb.Param[k] = v
				}
				saveNode++
				code := sp.SpExec(Sp, wb)
				if code != 200 {
					errorNode++
					seelog.Error("~~", Sp, wb)
				}
			}
			//else {
			//	if len(val) > 1 {
			//		fd.save(ele, val)
			//	}
			//}
		}
		if saveNode == 0 {
			seelog.Error("没在匹配到数据.出错正则表达式:", fd.ErrorReg, "\n\t", param.Param)
		} else {
			seelog.Info("save node=", saveNode, ";spit length=", len(fd.items), ";最后没有匹配的正则表达式", fd.ErrorReg)
		}
	}
}

func (fd *FmtData) getParam(html string, param *utils.Element) (map[string]interface{}) {
	res := make(map[string]interface{})
	ns := param.AllNodes()
	for _, n := range ns {
		if n.Name() == "Ref" {
			n = plugins.ElementMap[n.MustAttr("Id")]
		}
		regTxt := strings.TrimSpace(n.Value)
		if len(regTxt) < 1 {
			seelog.Error("没有设置正则表达式", n.ToString())
			return res
		}
		re := regexp.MustCompile(regTxt)
		gs := re.FindAllStringSubmatch(html, -1)
		if len(gs) >= 1 {
			Name := n.MustAttr("Name")
			Index := n.MustAttr("Index")
			ns := strings.Split(Name, ",")
			is := strings.Split(Index, ",")
			Sprintf, sfb := n.AttrValue("Sprintf")
			if len(ns) == len(is) && len(gs[0]) > len(ns) {
				for i := 0; i < len(ns); i++ {
					ni, _ := strconv.Atoi(is[i])
					if ni < 0 {
						ni = len(gs[0]) + ni
					}
					if len(gs[0]) > ni && ni > 0 {
						if sfb {
							res[ns[i]] = fmt.Sprintf(Sprintf, gs[0][ni])
						} else {
							res[ns[i]] = gs[0][ni]
						}
					} else {
						return nil
					}
				}
			} else {
				//seelog.Warn("配置错误，长度不匹配:", Name, " ~~ ", Index)
				return nil
			}
		} else {
			fd.ErrorReg = regTxt
			return nil
		}
	}
	return res
}

func (fd *FmtData) save(ele *utils.Element, val map[string]string) {
	SqlEle := ele.Node("Sql")
	if SqlEle != nil {
		Check := SqlEle.Node("Check")
		Insert := SqlEle.Node("Insert")
		Update := SqlEle.Node("Update")
		ip, ext := fd.sqlParam(val, Check.MustAttr("Param"))
		if !ext {
			return
		}
		c, cExt, _ := fd.Dao.QueryInt(Check.Value, ip...)
		if cExt && c >= 1 {
			if Update != nil {
				bp, ext := fd.sqlParam(val, Update.MustAttr("Param"))
				if ext {
					fd.Dao.Exec(Update.Value, bp...)
				}
			}
		} else {
			if Insert != nil {
				ip, ext := fd.sqlParam(val, Insert.MustAttr("Param"))
				if ext {
					fd.Dao.Exec(Insert.Value, ip...)
				}
			}
		}
	}
	seelog.Info(";", val)
}

func (fd *FmtData) sqlParam(val map[string]string, nd string) ([]interface{}, bool) {
	var p []interface{}
	ns := strings.Split(nd, ",")
	for _, n := range ns {
		v, ve := val[n]
		if ve {
			p = append(p, v)
		} else {
			return p, false
		}
	}
	return p, true
}

func (fd *FmtData) Call(ele *utils.Element, html string) {
	var rfd = reflect.ValueOf(fd)
	for _, n := range ele.AllNodes() {
		inputs := make([]reflect.Value, 2)
		inputs[0] = reflect.ValueOf(n)
		inputs[1] = reflect.ValueOf(html)
		m := rfd.MethodByName(n.Name())
		m.Call(inputs)
	}
}
