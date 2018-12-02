package cron

import (
	"github.com/ecdiy/goserver/plugins"
	"github.com/ecdiy/goserver/utils"
	"reflect"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/plugins/core"
)

var AppCron = New()

func init() {
	plugins.Plugins["Cron"] = func(ele *utils.Element) {
		job := &core.WebExec{Ele: ele}
		job.WebExec = reflect.ValueOf(job)
		spec, spb := ele.AttrValue("Spec")
		if spb && len(spec) > 1 {
			seelog.Info("Add Job:", spec)
			AppCron.AddFunc(spec, job.Job)
			in := ele.Attr("Init", "0")
			if in == "1" {
				job.Job()
			}
			if !AppCron.running {
				AppCron.Start()
			}
		} else {
			job.Job()
		}
	}
}
