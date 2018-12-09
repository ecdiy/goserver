package sp

/**
存储过程参数规则
1.ua
2.in*
3.wx*
4.gin*
5.wk*
 */

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins"
)

//--
func ParamUa(wb *utils.Param, p *SpParam) (interface{}, int) {
	return wb.Ua, 200
}

func ParamIn(ctx *utils.Param, p *SpParam) (interface{}, int) {
	v := ctx.String(p.ParamName)
	if v == "" {
		v = p.DefaultVal
	}
	return v, 200
}



//-----------

func ginWk(verify plugins.BaseFun, unFindCode int) func(wb *utils.Param, p *SpParam) (interface{}, int) {

	return func(wb *utils.Param, p *SpParam) (interface{}, int) {
		v2, b2 := wb.Context.Get(p.ParamName)
		if b2 {
			return v2, 200
		}
		verify(wb)
		v2, b2 = wb.Context.Get(p.ParamName)
		if b2 {
			return v2, 200
		}
		if unFindCode == 401 {
			return 0, 401
		} else {
			return p.DefaultVal, 200
		}
	}
}
