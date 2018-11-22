package core

import (
	"github.com/ecdiy/goserver/utils"
	"reflect"
	"github.com/cihub/seelog"
)

func (app *Module) Cron(ele *utils.Element) {

	job := &WebExec{ele: ele}
	job.webExec = reflect.ValueOf(job)
	spec, spb := ele.AttrValue("Spec")
	if spb && len(spec) > 1 {
		seelog.Info("Add Job:", spec)
		AppCron.AddFunc(spec, job.job)
		in := ele.Attr("Init", "0")
		if in == "1" {
			job.job()
		}
	} else {
		job.job()
	}

}
