package main

import (
	"goserver/utils"
	"reflect"
	"strings"
	"github.com/cihub/seelog"
	"regexp"
	"goserver/gpa"
	"strconv"
)

type fmtData struct {
	items []string
	dao   *gpa.Gpa
}

func (fd *fmtData) Spit(ele *utils.Element, html string) {
	regTxt := strings.TrimSpace(ele.Node("Regexp").Value)
	if len(regTxt) < 1 {
		seelog.Error("没有设置正则表达式", ele.ToString())
		return
	}
	Spit, SpitExt := ele.AttrValue("SpitString")
	if SpitExt {
		fd.items = strings.Split(html, Spit)
	}
	if fd.items != nil {
		ItemInclude, ItemIncludeExt := ele.AttrValue("ItemInclude")
		re := regexp.MustCompile(regTxt)
		for _, it := range fd.items {
			if ItemIncludeExt {
				if strings.Index(it, ItemInclude) < 0 {
					continue
				}
			}
			gs := re.FindAllStringSubmatch(it, -1)
			if len(gs) >= 1 {
				fd.save(ele, gs[0])
			}
		}
	}
	seelog.Info("items count:", len(fd.items))
}

func (fd *fmtData) save(ele *utils.Element, bd []string) {
	SqlEle := ele.Node("Sql")
	if SqlEle != nil {
		Check := SqlEle.Node("Check")
		Insert := SqlEle.Node("Insert")
		Update := SqlEle.Node("Update")
		ip, ext := fd.sqlParam(bd, Check.MustAttr("Param"))
		if !ext {
			return
		}
		c, cExt, _ := fd.dao.QueryInt(Check.Value, ip...)
		if cExt && c >= 1 {
			if Update != nil {
				bp, ext := fd.sqlParam(bd, Update.MustAttr("Param"))
				if ext {
					fd.dao.Exec(Update.Value, bp...)
				}
			}
		} else {
			if Insert != nil {
				ip, ext := fd.sqlParam(bd, Insert.MustAttr("Param"))
				if ext {
					fd.dao.Exec(Insert.Value, ip...)
				}
			}
		}
	}
	seelog.Info(";", bd)
}

func (fd *fmtData) sqlParam(param []string, nd string) ([]interface{}, bool) {
	var p []interface{}
	ns := strings.Split(nd, ",")
	for _, n := range ns {
		ni, _ := strconv.Atoi(n)
		if ni < 0 {
			ni = len(param) + ni
		}
		if len(param) > ni && ni >= 0 {
			p = append(p, param[ni])
		} else {
			return p, false
		}
	}
	return p, true
}

func (fd *fmtData) call(ele *utils.Element, html string) {
	var rfd = reflect.ValueOf(fd)
	for _, n := range ele.AllNodes() {
		inputs := make([]reflect.Value, 2)
		inputs[0] = reflect.ValueOf(n)
		inputs[1] = reflect.ValueOf(html)
		m := rfd.MethodByName(n.Name())
		m.Call(inputs)
	}
}
