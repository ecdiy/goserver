package demos

import (
	"testing"
	"github.com/ecdiy/goserver/utils"
	"fmt"
	"regexp"
)

func Test_Param(t *testing.T) {
	utils.EnvParamSet("a", "123")
	utils.EnvParamSet("xx", "1x23")

	re := regexp.MustCompile("[$][{]([^}]*)[}]")
	gs := re.FindAllStringSubmatch("${a}/${xx}", -1)

	fmt.Println(gs)

	fmt.Println(utils.FmtVal("${a}"))
	fmt.Println(utils.FmtVal("${a}/${bx}/fa"))
}
