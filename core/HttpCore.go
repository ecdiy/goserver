package core

import (
	"github.com/ecdiy/goserver/utils"
	"encoding/json"
	"github.com/cihub/seelog"
	"strings"
	"github.com/ecdiy/goserver/webs"
	"errors"
)

type HttpCore struct {
	html string
}

func (we *HttpCore) DoHttp(ele *utils.Element, param *webs.Param) error {
	GetUrl, gb := ele.AttrValue("GetUrl")
	if gb {
		http := utils.Http{}
		qp, qpExt := ele.AttrValue("QueryParam")
		if qpExt {
			qps := strings.Split(qp, ",")
			for _, p := range qps {
				GetUrl = strings.Replace(GetUrl, "${"+p+"}", param.String(p), -1)
			}
		}
		html, e := http.Get(GetUrl)
		if e == nil {
			return we.parseHtml(html, ele, param)
		} else {
			seelog.Error("HttpGet Fail:", GetUrl)
			return errors.New("HTTP get fail.")
		}
	}
	return nil
}

func (we *HttpCore) parseHtml(html string, ele *utils.Element, param *webs.Param) error {
	hs := ele.AllNodes()
	if hs != nil {
		we.html = html
		for _, h := range hs {
			var err error
			if h.Name() == "Map" {
				err = we.doMap(h, param)
			} else if h.Name() == "Json" {
				err = we.doJson(h, param)
			} else if h.Name() == "FmtData" {
				err = we.doFmtData(h, param)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (we *HttpCore) doFmtData(ele *utils.Element, param *webs.Param) error {
	//TODO
	return nil
}

func (we *HttpCore) doMap(ele *utils.Element, param *webs.Param) error {
	ns := ele.AllNodes()
	if ns != nil {
		for _, n := range ns {
			if n.Name() == "Body" {
				param.Param [n.MustAttr("Name")] = we.html
			} else {
				seelog.Error("没有实现的Map节点：" + n.Name())
				return errors.New("没有实现的Map节点:" + n.Name())
			}
		}
		return nil
	} else {
		seelog.Error("no map node.")
		return nil
	}
}

func (we *HttpCore) doJson(ele *utils.Element, param *webs.Param) error {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(we.html), &js)
	if err == nil {
		ns := ele.AllNodes()
		if ns != nil {
			for _, n := range ns {
				if n.Name() == "Put" {
					jk := n.MustAttr("JsonKey")
					v, vExt := js[jk]
					if !vExt {
						seelog.Error("JSON KEY NOT FIND:"+jk, we.html)
						return errors.New("JSON KEY NOT FIND:" + jk)
					}
					param.Param[n.MustAttr("Name")] = v
				}
			}
		}
		return nil
	} else {
		seelog.Error("DoJson fail.", we.html)
		return err
	}
}
