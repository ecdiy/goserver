# goserver
类似JAVA的spring.

#两种方式使用
*  配置文件，类似spring的配置文件，一个xml节点一个功能
*  建一个go main导入插件方式
*  cmd/App.go 默认导入了所有实现的插件，可以根据项目的需要，只导入自己需要的

## 如何注册插件
* web 插件 
```
import (
	"github.com/ecdiy/goserver/utils"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/plugins/web"
)
.....
web.RegisterWeb("CaptchaNew", func(xml *utils.Element) func(wb *utils.Param) {
		return func(wb *utils.Param) {
			wb.OK(captcha.New())
		}
})
```
* 插件
```
import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/gpa"
)

func init() {
	plugins.RegisterPlugin("Verify", func(ele *utils.Element) interface{} {
	...
	})
```


###  goserver 功能介绍
* 存储过程映射成JSON接口
* 权限验证
* 认证码 
* 定时任务（执行SQL，爬虫） 
* 模版
* 文件上传
* Web服务,静态资源
* 图片缩放
* Lua脚本支持
* WebSocket
* 二维码
* .........

更多的项目文档请参考：

https://www.itgeek.top/p/goserver

```angular2html
goserver 配置文件
``` 
 

# 使用需要什么技术
 sql
* 项目为go项目，你可以不懂go,下载对应平台的二进制包即可

# 本项目适合人群
* 想快速开发
* 前后端分离 
* 网站开发 (nuxt.js/vue+json数据请求接口)
* JSON接口  (APP开发)
* 大量的SQL操作
* 定时任务+爬虫 格式化数据 示例: https://www.itgeek.top/p/goserver/27 

# 示例说明
* 导入SQL： demos\goserver.sql
* vue : 
  npm install  
  npm run serve
  
* 配置 nginx  demos\demo\vue\admin.conf
  
* goserver goserver-dev.xml
  
# 谁在使用     
* ITGeek.top 所有后台为5个goserver实例，前台为nuxt.js。
* 加QQ群671735112，提交你的作品，讨论GoServer新功能 
 ....
    