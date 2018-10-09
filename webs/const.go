package webs

import (
	"regexp"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"utils/gpa"
	"utils"
	"strings"
)

const (
	DefaultPageSize = 20
	SqlSpAll  = "select name,CONVERT(param_list USING utf8) param_list,`comment` from mysql.proc c where db=DATABASE() and `type`='PROCEDURE'"
	SqlSpInfo = "select name,CONVERT(param_list USING utf8) param_list,`comment` from mysql.proc c where db=DATABASE() and `type`='PROCEDURE' and name=?"
)

var (
	Gin         = gin.New()
	Gpa         *gpa.Gpa

	spCache     = make(map[string]*Sp)
	spReloadFun = make(map[string]func(c *Param) *UserBase)
	UaH5        *regexp.Regexp
)

func init() {
	var e error
	UaH5, e = regexp.Compile(`(MIDP)|(WAP)|(UP.Browser)|(Smartphone)|(Obigo)|(Mobile)|(AU.Browser)|(wxd.Mms)|(WxdB.Browser)|(CLDC)|(UP.Link)|(KM.Browser)|(UCWEB)|(SEMC\-Browser)|(Mini)|(Symbian)|(Palm)|(Nokia)|(Panasonic)|(MOT\-)|(SonyEricsson)|(NEC\-)|(Alcatel)|(Ericsson)|(BENQ)|(BenQ)|(Amoisonic)|(Amoi\-)|(Capitel)|(PHILIPS)|(SAMSUNG)|(Lenovo)|(Mitsu)|(Motorola)|(SHARP)|(WAPPER)|(LG\-)|(LG/)|(EG900)|(CECT)|(Compal)|(kejian)|(Bird)|(BIRD)|(G900/V1.0)|(Arima)|(CTL)|(TDG)|(Daxian)|(DAXIAN)|(DBTEL)|(Eastcom)|(EASTCOM)|(PANTECH)|(Dopod)|(Haier)|(HAIER)|(KONKA)|(KEJIAN)|(LENOVO)|(Soutec)|(SOUTEC)|(SAGEM)|(SEC\-)|(SED\-)|(EMOL\-)|(INNO55)|(ZTE)|(iPhone)|(Android)|(Windows CE)|(Wget)|(Java)|(curl)|(Opera)/`)
	if e != nil {
		seelog.Error("h5 ua error.", e)
	}
}

func Init(db string, models ...interface{}) {
	if utils.EnvIsDev {
		ip := utils.GetIp()
		utils.EnvParamSet("ImgHost", "http://"+ip)
	}

	if strings.Index(db, ":") < 0 {
		utils.EnvParamSet("DbDsn", "root:root@tcp(127.0.0.1:3306)/" + db+
			"?timeout=30s&charset=utf8mb4&parseTime=true")
	} else {
		utils.EnvParamSet("DbDsn", db)
	}
	dsn := utils.EnvParam("DbDsn")
	Gpa = gpa.Init(utils.EnvParam("DbDriver"), dsn, models...)
	seelog.Info(dsn)

}
