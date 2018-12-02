package http

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/webs"
	"strings"
	"github.com/cihub/seelog"
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
		seelog.Info("GetUrl:", GetUrl)
		html, e := http.Get(GetUrl)
		if e == nil {
			return we.parseHtml(html, ele, param)
		} else {
			seelog.Error("HttpGet Fail:", GetUrl)
			return errors.New("HTTP get fail.")
		}
	} else {
		//TODO　POST
		seelog.Warn("没有GetUrl，POST?")
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
				err = we.parseMap(h, param)
			} else if h.Name() == "Json" {
				err = we.parseJson(h, param)
			} else if h.Name() == "List" {
				err = we.parseList(h, param)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
