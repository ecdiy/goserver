package http

import (
	"strings"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/webs"
	"github.com/ecdiy/goserver/plugins"
)

func (we *HCore) parseList(ele *utils.Element, param *webs.Param) error {
	Begin := ele.Node("Begin")
	html := we.html
	if Begin != nil {
		ind := strings.Index(html, Begin.Value)
		if ind < 0 {
			seelog.Error("内容没有开始（Begin）标记", Begin.Value)
			return nil
		}
		html = html[ind:]
	}
	End := ele.Node("End")
	if End != nil {
		ind := strings.Index(html, End.Value)
		if ind < 0 {
			seelog.Error("内容没有开始（End）标记:", End.Value, html)
			return nil
		}
	}
	fd := &FmtData{Dao: plugins.GetGpa(ele)}
	fd.Spit(ele, html, param)
	return nil
}
