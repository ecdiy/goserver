package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/cihub/seelog"
	"os"
	"strings"
	"os/exec"
	"syscall"
	"io"
)

func (app *Module) Upgrade(ele *utils.Element) {
	up := &UpgradeImpl{CheckUrl: ele.MustAttr("CheckUrl"),
		UpgradeFlag: ele.Attr("UpgradeFlag", "upgrade")}
	up.Bin = os.Args[0]
	idx := strings.LastIndex(os.Args[0], string(os.PathSeparator))
	if idx > 0 {
		up.Bin = os.Args[0][idx+1:]
	}
	if strings.Index(up.Bin, up.UpgradeFlag) == 0 {
		up.BinReplaceCheck()
	} else {
		up.WorkPath, _ = os.Getwd()
		spec, spb := ele.AttrValue("Spec")
		if spb && len(spec) > 1 {
			seelog.Info("Add upgrade Job:", spec)
			up.DoUpgrade()
			AppCron.AddFunc(spec, up.DoUpgrade)
		} else {
			up.DoUpgrade()
		}
	}
}

type UpgradeImpl struct {
	Bin, UpgradeFlag, CheckUrl, BinVersion, WorkPath string
}

func (u *UpgradeImpl) BinReplaceCheck() {
	initBin := u.WorkPath + u.Bin[len(u.UpgradeFlag):]
	seelog.Info(initBin, "~~", os.Args[0])
	CopyFile(os.Args[0], u.WorkPath+u.Bin[len(u.UpgradeFlag):])
	u.RunCmd(initBin)
	os.Exit(1)
}

func (u *UpgradeImpl) RunCmd(newBin string) {
	cmd := exec.Command(newBin)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Start()
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func (u *UpgradeImpl) DoUpgrade() {
	http := utils.Http{}
	val, e := http.Json("GET", u.CheckUrl)
	if e != nil {
		seelog.Error("升级查检有错:", u.CheckUrl, "\n", e)
		return
	}
	bv, bvExt := val["BinVersion"]
	cv := utils.EnvParam("BinVersion")
	if bvExt && cv != bv {
		seelog.Info("有Bin新版本:", bv, ";当前版本=", cv, val["BinDesc"])
		u.BinUpgrade(val["BinUrl"].(string))
	}
	ui, uiExt := val["StaticVersion"].(string)
	if uiExt && utils.EnvParam("StaticVersion") != ui {
		u.StaticUpdate(val["StaticUrl"].(string))
	}
}

func (u *UpgradeImpl) BinUpgrade(UrlBinZip string) {
	http := utils.Http{}
	http.GetUnzip(UrlBinZip, u.WorkPath+"/"+u.UpgradeFlag)
	upgrade := u.WorkPath + "/" + u.UpgradeFlag + "/" + u.Bin
	newBin := u.WorkPath + "/" + u.UpgradeFlag + "_" + u.Bin
	seelog.Info(os.Args[0], " ~ ", upgrade, " ~ ", newBin)
	os.Rename(upgrade, newBin)
	u.RunCmd(newBin)
	os.Exit(1)
}

func (u *UpgradeImpl) StaticUpdate(UrlStaticZip string) {
	http := utils.Http{}
	http.GetUnzip(UrlStaticZip, u.WorkPath)
}
