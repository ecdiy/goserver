package webs

import "github.com/gin-gonic/gin"

func GetUa(ctx *gin.Context) string {
	ua := ctx.Request.UserAgent()
	if len(ua) == 0 {
		return "web"
	}
	if UaH5.MatchString(ua) {
		return "h5"
	}
	if UaSeo.MatchString(ua) {
		return "seo"
	}
	return "web"
}
