
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
   
## 下一步开发计划 DevOps
* DevOps 基于GoServer应用项目,汇集：笔记，任务管理，软件库
* 笔记记录内容到本机
* 任务管理（可视化配置）
* 开发软件安装