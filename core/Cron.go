package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/utils/cron"
	"reflect"
	"github.com/cihub/seelog"
)

func (app *Module) Cron(ele *utils.Element) {
	cron := cron.New()
	job := &WebExec{ele: ele}
	job.webExec = reflect.ValueOf(job)
	spec, spb := ele.AttrValue("Spec")
	if spb && len(spec) > 1 {
		seelog.Info("Add Job:", spec)
		cron.AddFunc(spec, job.job)
		in := ele.Attr("Init", "0")
		if in == "1" {
			job.job()
		}
	} else {
		job.job()
	}
	cron.Start()
}
