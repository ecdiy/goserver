package core

import (
	"reflect"
	"github.com/cihub/seelog"
	"os"
	"strings"
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/utils/cron"
	"github.com/ecdiy/goserver/plugins"
)

//---
var app = reflect.ValueOf(new(Module))

var initAfterFun []func()               //xml 分析完后的回调函数
var AppCron *cron.Cron

func StartCore() {
	seelog.Info("version: 0.3")
	defer func() {
		if err := recover(); err != nil {
			seelog.Flush()
			seelog.Error("un catch error;", err)
			panic(err)
		}
	}()
	if len(os.Args) < 2 {
		seelog.Error("没有指定配置文件")
		return
	}
	if strings.Index(os.Args[1], "-dev.xml") > 0 {
		utils.EnvParamSet("profile", utils.EnvDev)
	} else {
		utils.EnvParamSet("profile", utils.EnvProd)
	}
	dom, err := utils.LoadByFile(os.Args[1])
	if err == nil {
		InvokeByXml(dom)
	} else {
		seelog.Error("读取配置文件出错:", os.Args[1], ",", err)
		return
	}
}

func InvokeByXml(ecXml *utils.Element) {
	AppCron = cron.New()
	ns := ""
	defer func() {
		seelog.Info("analysis element:", ns, "\n\tdata:",plugins. Data)
	}()
	allNode := ecXml.AllNodes()
	seelog.Info("配置节点数:", len(allNode), ",WebPlugin:", plugins.WebPlugins, ",Plugin:", plugins.Plugins)
	for _, n := range allNode {
		ns += n.Name() + ";"
		//if IsReload {
		//	canReload, canReloadBool := n.AttrValue("canReload")
		//	if !canReloadBool || canReload != "1" {
		//		continue
		//	}
		//	log.Info("~~~reload~~~", n.Name())
		//}
		w, we := plugins.WebPlugins[n.Name()]
		if we {
			mtd := strings.ToUpper(n.Attr("Method", "Get"))
			if strings.Index(mtd, "GET") >= 0 {
				plugins.GetGin(n).GET(n.MustAttr("Url"), w(n))
			}
			continue
		}

		p, pFd := plugins.Plugins[n.Name()]
		if pFd {
			p(n)
			continue
		}

		inputs := make([]reflect.Value, 1)
		inputs[0] = reflect.ValueOf(n)
		m := app.MethodByName(n.Name())
		if m.IsValid() {
			m.Call(inputs)
		} else {
			panic("没有实现的方法:" + n.Name())
		}

	}

	if len(initAfterFun) > 0 {
		for _, fun := range initAfterFun {
			fun()
		}
	}
	AppCron.Start()
}
