package core

import (
	"github.com/ecdiy/goserver/utils"
	"encoding/json"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/webs"
	"errors"
)




func (we *HttpCore) parseMap(ele *utils.Element, param *webs.Param) error {
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

func (we *HttpCore)parseJson(ele *utils.Element, param *webs.Param) error {
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
