
###  项目介绍
goserver 是一个组件式后台服务，把常用的功能用配置的方式加载，如数据库CRUD，认证码 ... 

```angular2html
goserver 配置文件
```

```angular2html

<Root>

    <Gpa DbDsn="root:root@tcp(127.0.0.1:3306)/wk-topic?timeout=30s&amp;charset=utf8mb4&amp;parseTime=true"/>

    <!-- 验证功能 -->
    <Verify Id="Verify" TplFunName="verify" GpaRef="Gpa" TokenName="token" RpcHost="127.0.0.1:7200"
            ResultFlagName="Login"/>

    <Web Port=":7002">
        <Static Url="/static" Path="${web.static}"/>
    </Web>

    <Sp WebRef="Web" GpaRef="Gpa">
        <ParamGin Prefix="gin" VerifyRef="Verify"/>
        <ParamWk Prefix="gk" VerifyRef="Verify"/>
        <!--
        1. RpcRef="Rpc"
        2. RpcHost="127.0.0.1:4200"
        3. Sql="select UserId from Token where Token=? and Ua=?"
        三种方式任选一种,未发现时返回 code=401. ParamGin 区别，未匹配时返回默认值0 ,code=200
         -->
        <!-- ReloadUrl 清空内存缓存 设置成 "" 为不能用URL访问的方法请求内存缓存 -->

        <Handle SpSuffix="Ajax" ReloadUrl="/spReload" Url="/sp/:sp"/>

        <!--<Template SpSuffix="Page" TemplatesDir="ui/view/geek/wk-topic/" LayoutDir="ui/view/geek/layout" Pages="/**/*,/*"-->
                  <!--LoginUrl="/user/login">-->
            <!--<extends>-->
                <!--<Item TplName="index-web">ui/view/geek/modules/topicList.html</Item>-->
                <!--<Item TplName="my/fav-web">ui/view/geek/modules/topicList.html</Item>-->
                <!--<Item TplName="my/follow-web">ui/view/geek/modules/topicList.html</Item>-->
            <!--</extends>-->
            <!--<map>-->
                <!--<Item TplName="detail" Url="/detail/:Id/:UserId"/>-->
                <!--<Item TplName="append" Url="/append/:TopicId"/>-->
                <!--<Item TplName="index" Url="/list/:CatId/:Page"/>-->
                <!--<Item TplName="index" Url="/page/:Page"/>-->
            <!--</map>-->
        <!--</Template>-->
    </Sp>

    <Upload WebRef="Web" Url="/api/Upload" TmpDir="./upload/temp/" DirUpload="${upload.dir}" ImgWidth="800" MainWidth="800"
            UrlPrefix="${upload.prefix}"/>
</Root>

```
### 插件功能描述

#### Gpa 数据库
```angular2html
   <Gpa DbDsn="root:root@tcp(127.0.0.1:3306)/wk-topic?timeout=30s&amp;charset=utf8mb4&amp;parseTime=true"/>
```

#### Verify 验证码功能
```angular2html
<Verify Id="Verify" TplFunName="verify" GpaRef="Gpa" TokenName="token" RpcHost="127.0.0.1:7200"
            ResultFlagName="Login"/>
```

#### Web  绑定端口，静态资源共享目录
```angular2html
  <Web Port=":7002">
        <Static Url="/static" Path="${web.static}"/>
  </Web>
```

#### Sp   数据库CRUD
```angular2html
    <Sp WebRef="Web" GpaRef="Gpa">
        <ParamGin Prefix="gin" VerifyRef="Verify"/>
        <ParamWk Prefix="gk" VerifyRef="Verify"/>
        <!--
        1. RpcRef="Rpc"
        2. RpcHost="127.0.0.1:4200"
        3. Sql="select UserId from Token where Token=? and Ua=?"
        三种方式任选一种,未发现时返回 code=401. ParamGin 区别，未匹配时返回默认值0 ,code=200
         -->
        <!-- ReloadUrl 清空内存缓存 设置成 "" 为不能用URL访问的方法请求内存缓存 -->

        <Handle SpSuffix="Ajax" ReloadUrl="/spReload" Url="/sp/:sp"/>
    </Sp>
```