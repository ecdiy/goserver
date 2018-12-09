package plugins

import (
	"github.com/ecdiy/goserver/utils"
	"reflect"
	"os"
	"strings"
)

func PutFunRun(fun func()) {
	if len(InitAfterFun) > 0 {
		go fun()
	} else {
		InitAfterFun = append(InitAfterFun, fun)
	}
}

func put(ele *utils.Element, v interface{}) {
	id, idb := ele.AttrValue("Id")
	if !idb {
		id = ele.Name()
	}
	_, de := Data[id]
	if de {
		panic("Id重复" + id)
	} else {
		Data[id] = v
	}
}

func getVerify(ele *utils.Element) BaseFun {
	VerifyId, vb := ele.AttrValue("VerifyRef")
	if !vb {
		VerifyId = "Verify"
	}
	return Data[VerifyId].(BaseFun)
}

func doSubElement(ele *utils.Element, obj interface{}) {
	ns := ele.AllNodes()
	if len(ns) > 0 {
		spv := reflect.ValueOf(obj)
		for _, e := range ns {
			inputs := make([]reflect.Value, 2)
			inputs[0] = reflect.ValueOf(e)
			inputs[1] = reflect.ValueOf(Data)
			m := spv.MethodByName(e.Name())
			m.Call(inputs)
		}
	}
}

func getFile(file string) string {
	dir := os.Args[1]
	fg := []string{"/", "\\"}
	for _, flg := range fg {
		lst := strings.LastIndex(dir, flg)
		if lst > 0 {
			return dir[0:lst+1] + file
		}
	}
	return file
}

