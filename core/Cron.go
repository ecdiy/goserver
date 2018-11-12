package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/xtools/cron"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/gpa"
	"reflect"

	"github.com/ecdiy/goserver/xtools"
)

func (app *Module) Cron(ele *utils.Element) {
	cron := cron.New()
	ns := ele.AllNodes()
	for _, n := range ns {
		job := &Job{job: n, dao: getGpa(n)}
		spec, spb := n.AttrValue("Spec")
		if spb {
			cron.AddFunc(spec, job.Run)
		} else {
			job.Run()
		}
	}
	cron.Start()
}

type Job struct {
	job *utils.Element
	dao *gpa.Gpa
}

func (job *Job) Run() {
	inputs := make([]reflect.Value, 0)
	m := reflect.ValueOf(job).MethodByName(job.job.Name())
	if m.IsValid() {
		m.Call(inputs)
	} else {
		seelog.Error("未知方法:" + job.job.Name())
	}
}

func (job *Job) Sql() {
	seelog.Info("~~Sql定时任务：", job.job.Value)
	job.dao.Exec(job.job.Value)
}
func (job *Job) Http() {
	getUrl, ext := job.job.AttrValue("GetUrl")
	seelog.Info("~~Http定时任务：", job.job.Value)
	if ext {
		http := utils.Http{}
		html, e := http.Get(getUrl)
		if e == nil {
			fd := &xtools.FmtData{Dao: job.dao}
			fmt := job.job.Node("Fmt")
			if fmt != nil {
				fd.Spit(job.job.Node("Fmt"), string(html))
			}
		} else {
			seelog.Error("HttpGet失败", getUrl)
		}
	}
}
