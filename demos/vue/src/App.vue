<template>
    <div id="app" style="background:#eee;">
        <card style="height: 50px;">
            <div class="content">
                <div class="fr">
                    <router-link to="/">首页</router-link>
                    <!--<router-link to="/user/login">登录</router-link>-->
                    <router-link to="/user/register">注册</router-link>
                    <a href="https://www.itgeek.top" target="_blank">ITGeek.top</a>
                    <a href="https://www.itgeek.top/p/goserver" target="_blank">文档</a>
                </div>
            </div>
        </card>
        <router-view></router-view>
        <BackTop :height="100" :bottom="200">
            <div class="top">返回顶端</div>
        </BackTop>
    </div>
</template>
<script>
    import Cookies from 'js-cookie';

    export default {
        data() {
            return {
                an: '', on: [], menu: []
            }
        },
        methods: {
            loginOut() {
                Cookies.remove(this.conf.cookieTokenName);
                this.$router.replace(this.conf.loginUrl)
            }
        }
    }
</script>
<style>
    html, body {
        width: 100%;
        height: 100%;
        background: #f0f0f0;
    }

    a:active, a:link, a:visited {
        color: #778087;
        text-decoration: none;
        word-break: break-all;
    }

    p {
        padding-top: 5px;
    }

    .app-main {
        width: 100%;
        height: 100%;
    }

    .fl {
        float: left;
    }

    .fr {
        float: right;
    }

    .small, .small a {
        padding-left: 10px;
        font-size: 11px;
        color: #646464;
        cursor: pointer;
    }

    .ft {
        margin-top: 10px;
        padding: 10px 0;
        background-color: #222;
        width: 100%;
        height: 140px;
        color: #fff;
        line-height: 30px;
    }

    .ft a {
        padding-left: 10px;
        color: #fff;
        line-height: 30px;
    }

    a {
        padding-left: 10px;
    }

    .content {
        min-width: 600px;
        max-width: 1100px;
        margin: 0 auto 0 auto;
    }

    .h li {
        float: left;
        padding-left: 40px;
        list-style: none;
        font-size: 18px;
    }

    .current {
        display: inline-block;
        font-size: 14px;
        line-height: 14px;
        padding: 5px 8px 5px 8px;
        margin-right: 5px;
        border-radius: 3px;
        background-color: #334;
        color: #fff;
    }

    a.current {
        color: #fff;
    }
</style>
