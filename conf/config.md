
sp :参数规则

```$xslt
gin:    取gin.Context 中Get方法的值，没有命中则返401
in:     各种参数的值，post,get,path ... 中取值
ua:     web,h5
```


```angular2html
<Root>
    <Parameter>
        <profile>dev</profile>
    </Parameter>

    <Gpa DbDsn="root:root@tcp(127.0.0.1:3306)/dinner-admin?timeout=30s&amp;charset=utf8mb4&amp;parseTime=true"/>

    <Web/>

    <Rpc Sql="select UserId from Token where Token=? and Ua=?" RpcHost="127.0.0.1:4200" GpaRef="Gpa"/>


    <!-- ReloadUrl 清空内存缓存 设置成 "" 为不能用URL访问的方法请求内存缓存 -->
    <WebSp SpSuffix="Ajax" ReloadUrl="/spReload" Url="/sp/:sp" WebRef="Web" GpaRef="Gpa">
        <ParamGin Prefix="gin" Method="GinDbParam" TokenName="token" RpcRef="Rpc"/>
    </WebSp>

    <!--<RpcSp  TokenName="token"-->
    <!--Url="/sp/:sp" Ext="Ajax"/>-->

    <!--<RpcHost RpcHost="127.0.0.1:4200" TokenName="token" Url="/spa/:sp" Ext="Admin"/>-->

    <!--<SpReload Url="/spReload"/>-->

    <WebRun Port=":4001" WebRef="Web"/>
</Root>

```

```angular2html
<WebSp SpSuffix="Ajax" ReloadUrl="/spReload" Url="/sp/:sp" WebRef="Web" GpaRef="Gpa">
    <ParamGin Prefix="gin" Method="GinDbParam" TokenName="token" RpcRef="Rpc"/>
</WebSp>

```