package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/utils/cron"
	"reflect"
)

func (app *Module) Cron(ele *utils.Element) {
	cron := cron.New()
	ns := ele.AllNodes()
	for _, n := range ns {
		job := &WebExec{ele: ele}
		job.webExec = reflect.ValueOf(job)
		spec, spb := n.AttrValue("Spec")
		if spb {
			cron.AddFunc(spec, job.job)
		} else {
			job.job()
		}
	}
	cron.Start()
}
