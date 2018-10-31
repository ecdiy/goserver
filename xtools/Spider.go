package xtools

import (
	"reflect"
	"strings"
	"github.com/cihub/seelog"
	"regexp"
	"goserver/gpa"
	"strconv"
	"goserver/utils"
)

type FmtData struct {
	items []string
	Dao   *gpa.Gpa
}

func (fd *FmtData) Spit(ele *utils.Element, html string) {
	Spit, SpitExt := ele.AttrValue("SpitString")
	if SpitExt {
		fd.items = strings.Split(html, Spit)
	}
	if fd.items != nil {
		ItemInclude, ItemIncludeExt := ele.AttrValue("ItemInclude")
		for _, it := range fd.items {
			if ItemIncludeExt {
				if strings.Index(it, ItemInclude) < 0 {
					continue
				}
			}
			val := fd.getParam(it, ele.Node("Param"))
			if len(val) > 1 {
				fd.save(ele, val)
			}
		}
	}
	if len(fd.items) == 1 {
		//seelog.Warn("?spider error?")
		//ioutil.WriteFile("", []byte(html), 066)
	}
	seelog.Info("items count:", len(fd.items))
}

func (fd *FmtData) getParam(html string, param *utils.Element) map[string]string {
	res := make(map[string]string)
	ns := param.AllNodes()
	for _, n := range ns {
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
			if len(ns) == len(is) && len(gs[0]) > len(ns) {
				for i := 0; i < len(ns); i++ {
					ni, _ := strconv.Atoi(is[i])
					if ni < 0 {
						ni = len(gs[0]) + ni
					}
					if len(gs[0]) > ni && ni > 0 {
						res[ns[i]] = gs[0][ni]
					}
				}
			}
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
