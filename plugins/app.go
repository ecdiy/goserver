package plugins

import (
	"reflect"
	"github.com/cihub/seelog"
	"os"
	"strings"
	"github.com/ecdiy/goserver/utils"
)

//---
var app = reflect.ValueOf(new(Module))

func StartCore() {
	seelog.Info("version: 0.1.1")
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
	ns := ""
	defer func() {
		seelog.Info("analysis element:", ns, "\n\tdata:", Data)
	}()
	allNode := ecXml.AllNodes()
	seelog.Info("配置节点数:", len(allNode), ",Plugin:", pluginsMap)
	for _, n := range allNode {
		ns += n.Name() + ";"
		p, pFd := pluginsMap[n.Name()]
		if pFd {
			pImpl := p(n)
			if pImpl != nil {
				put(n, pImpl)
			}
		} else {
			inputs := make([]reflect.Value, 1)
			inputs[0] = reflect.ValueOf(n)
			m := app.MethodByName(n.Name())
			if m.IsValid() {
				m.Call(inputs)
			} else {
				seelog.Flush()
				panic("没有实现的方法:" + n.Name())
			}
		}
	}

	if len(InitAfterFun) > 0 {
		for _, fun := range InitAfterFun {
			fun()
		}
	}

}

