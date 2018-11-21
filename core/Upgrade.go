package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/utils/overseer"
	"github.com/ecdiy/goserver/utils/overseer/fetcher"
)

func (app *Module) Upgrade(ele *utils.Element) {
	up := &UpgradeImpl{CheckUrl: ele.MustAttr("CheckUrl")}
	spec, spb := ele.AttrValue("Spec")
	if spb && len(spec) > 1 {
		seelog.Info("Add upgrade Job:", spec)
		up.DoUpgrade()
		AppCron.AddFunc(spec, up.DoUpgrade)
	} else {
		up.DoUpgrade()
	}
}

type UpgradeImpl struct {
	CheckUrl, BinVersion, WorkPath string
}

func (u *UpgradeImpl) DoUpgrade() {
	//	vJson, e := ioutil.ReadFile(uiVerPath)
	http := utils.Http{}
	val, e := http.Json("GET", u.CheckUrl)
	if e != nil {
		seelog.Error("升级查检有错:", u.CheckUrl, "\n", e)
		return
	}
	bv, bvExt := val["BinVersion"]
	if bvExt && utils.EnvParam("BinVersion") != bv {
		seelog.Info("有Bin新版本:", bv, val["BinDesc"])
		u.BinUpgrade(val["BinUrl"].(string))
	}

	ui, uiExt := val["StaticVersion"].(string)
	if uiExt && utils.EnvParam("StaticVersion") != ui {
		u.StaticUpdate(val["StaticUrl"].(string))
	}
	//seelog.Info("有UI新版本:", uiVersion, s["staticDesc"])
	//StaticUpdate()
}

func (u *UpgradeImpl) BinUpgrade(UrlBinZip string) {
	overseer.Run(overseer.Config{
		Program: func(state overseer.State) {},
		Fetcher: &fetcher.HTTP{
			URL: UrlBinZip,
			//e.g.http://localhost:4000/binaries/app-linux-amd64
		},
	})
	//http := utils.Http{}
	//http.GetUnzip(UrlBinZip, u.WorkPath+"/upgrade")
	//bin :=u.WorkPath  + "/upgrade/assistant.exe"
	//b, e := lang.FileExists(bin)
	//if b && e == nil {
	//	cmd := exec.Command(bin, "kill="+strconv.Itoa(os.Getpid()))
	//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//	cmd.Stdout = os.Stdout
	//	cmd.Start()
	//}
}

//func SysUpgrade() bool {
//	if len(os.Args) > 1 {
//		kill := tools.GetOsParam("kill")
//		if len(kill) > 0 {
//			//--升级
//			pid, _ := strconv.ParseInt(kill, 10, 0)
//			pro, err := os.FindProcess(int(pid))
//			seelog.Info("..........kill old pid=", pid)
//			if err == nil {
//				pro.Kill()
//				pro.Release()
//			}
//		}
//
//		if strings.Index(os.Args[0], "_new") > 0 {
//			exe := strings.Replace(os.Args[0], "_new", "", -1)
//			for i := 0; i < 10; i++ {
//				time.Sleep(2 * time.Second)
//				delErr := os.Remove(exe)
//				if delErr != nil || utils.IsExist(exe) {
//					seelog.Error(".....升级主程序失败，旧程序没有删除 :"+exe, delErr)
//				} else {
//					seelog.Info(".....升级主程序[", os.Getpid(), "]替换旧程序成功!....... ")
//					utils.CopyFile(os.Args[0], exe)
//					tools.RunRestartCmd(exe)
//					break
//				}
//			}
//			return true
//		} else {
//			appNewFile := os.Args[0]
//			idx := strings.LastIndex(appNewFile, ".")
//			if idx > 0 {
//				appNewFile = appNewFile[0:idx] + "_new" + appNewFile[idx:]
//			} else {
//				appNewFile += "_new"
//			}
//			seelog.Info(".....升级主程序,清理临时文件:" + appNewFile)
//			if utils.IsExist(appNewFile) {
//				os.Remove(appNewFile)
//			}
//		}
//	}
//	return false
//}

func (u *UpgradeImpl) StaticUpdate(UrlStaticZip string) {
	http := utils.Http{}
	http.GetUnzip(UrlStaticZip, u.WorkPath)
}
