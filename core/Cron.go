package core

import (
	"goserver/utils"
	"goserver/xtools/cron"
)

//cron := cron.New()
//cron.AddFunc("0 * * * * ?", func() {

func (app *Module) Cron(ele *utils.Element) {
	cron := cron.New()
	cron.AddFunc(ele.MustAttr("Spec"), func() {

	})
	cron.Start()
	defer cron.Stop()
}
