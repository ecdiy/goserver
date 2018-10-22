package main

import (
	"testing"
	"io/ioutil"
	"fmt"
	"goserver/utils"
	"goserver/gpa"
)

// https://github.com/search?l=&p=2&q=stars%3A%3E1000&ref=advsearch&type=Repositories&utf8=%E2%9C%93
/*
1.spit 分隔成数组.
*/

func Test_Github(t *testing.T) {
	bs, err := ioutil.ReadFile("c:/gopath/src/goserver/cmd/testdata/github_search_result.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	dom, err := utils.LoadByFile("c:/gopath/src/goserver/cmd/testdata/github_search_result.xml")
	FmtData(dom, string(bs))
}

func FmtData(ele *utils.Element, html string) {
	fd := &fmtData{dao: gpa.InitGpa("root:root@tcp(127.0.0.1:3306)/wk-site?timeout=30s&charset=utf8mb4&parseTime=true")}
	fd.call(ele, html)
}
