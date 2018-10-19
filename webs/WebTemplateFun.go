package webs

import (
	"html/template"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"
)

var FunConstantMaps = template.FuncMap{
	"param": func(name string) string { return "" }, //占位
	"set": func(m map[string]interface{}, n string, v interface{}) string {
		m[n] = v
		return ""
	},
	"get": func(m map[string]interface{}, n string) interface{} {
		return m[n]
	},
	"concat": func(s ...string) string {
		ss := ""
		for _, si := range s {
			ss += si
		}
		return ss
	},
	"unescaped": func(x string) template.HTML {
		return template.HTML(x)
	},
	"add": func(v ... int) int {

		vs := 0
		for _, vv := range v {
			vs += vv
		}
		return vs
	},
	"marshal": func(in interface{}) template.HTML {
		jsonStr, _ := json.Marshal(in)
		return template.HTML(string(jsonStr))
	},
	"neq": func(a, b interface{}) bool {
		return fmt.Sprint(a) != fmt.Sprint(b)
	},
	"eq": func(a, b interface{}) bool {
		r := a == b
		if r {
			return r
		} else {
			return fmt.Sprint(a) == fmt.Sprint(b)
		}
	},
	"gt": func(a, b interface{}) bool {
		ai, _ := strconv.ParseInt(fmt.Sprint(a), 10, 0)
		bi, _ := strconv.ParseInt(fmt.Sprint(b), 10, 0)
		return ai > bi
	},
	"page": func(url string, countF interface{}, p, current int) template.HTML {
		count, _ := strconv.Atoi(fmt.Sprint(countF))
		pall := count / p
		if pall*p < count {
			pall = pall + 1
		}
		if pall <= 1 {
			return template.HTML("")
		}
		h := `<div style="clear: both; padding-top: 10px;"><ul class="ivu-page"><span class="ivu-page-total">共 ` +
			fmt.Sprint(count) + ` 条</span>`
		if current > 1 {
			r := strings.Replace(url, "{}", fmt.Sprint(current-1), -1)
			h += `<li title="上一页" class="ivu-page-prev"><a href="` + r + `"><i class="ivu-icon ivu-icon-ios-arrow-back"></i></a></li>`
		}
		bg := current - 3
		if bg < 1 {
			bg = 1
		}
		end := current + 3
		if end > pall {
			end = pall
		}
		for i := bg; i < end; i++ {
			is := fmt.Sprint(i)
			r := strings.Replace(url, "{}", is, -1)
			c := "ivu-page-item"
			if i == current {
				c += " ivu-page-item-active"
			}
			h += `<li class="` + c + `"><a href="` + r + `">` + is + `</a></li>`
		}
		if current < pall {
			r := strings.Replace(url, "{}", fmt.Sprint(current+1), -1)
			h += `<li title="下一页" class="ivu-page-next ivu-page-disabled"><a href="` + r + `"><i
                class="ivu-icon ivu-icon-ios-arrow-forward"></i></a></li>`
		}
		h += `</ul></div>`
		return template.HTML(h)
	},
}
