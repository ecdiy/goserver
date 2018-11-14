<template>
    <card>
        <div style="width:300px;margin: auto">
            <div class="login-con">
                <Card :bordered="false">
                    <p slot="title">
                        <Icon type="ios-log-in"></Icon>
                        欢迎登录
                    </p>
                    <div class="form-con">
                        <Form :model="form" ref="form">
                            <FormItem prop="Username" :rules="rules.Username" label="用户名" :error="err.Username">
                                <Input v-model="form.Username" placeholder="请输入用户名">
                                <span slot="prepend">
                                    <Icon :size="16" type="ios-person"></Icon>
                                </span>
                                </Input>
                            </FormItem>
                            <FormItem prop="Password" :rules="rules.Password" label="密码" :error="err.Password">
                                <Input type="password" v-model="form.Password" @on-enter="handleSubmit" placeholder="请输入密码">
                                <span slot="prepend">
                                    <Icon :size="14" type="locked"></Icon>
                                </span>
                                </Input>
                            </FormItem>
                            <FormItem prop="password" v-show="authImg!=''" :error="err.Captcha">
                                <Input type="text" v-model="form.Digits" placeholder="请输入认证码">
                                </Input>
                                <img :src="authImg" @click="loadCaptcha"/>
                            </FormItem>
                            <FormItem>
                                <Button @click="handleSubmit" type="primary" long>登录</Button>
                            </FormItem>
                        </Form>
                    </div>
                </Card>
            </div>
        </div>
    </card>
</template>


<script>
    import Cookies from 'js-cookie';

    export default {
        props: ['prefix', 'goUrl'],
        data() {
            return {
                reqUrl: '/sp/UserLogin', authImgUrl: '/Captcha?t=', captchaNew: '/CaptchaNew',
                authImg: "", form: {Captcha: "", Username: '', Digits: ''},
                err: {Username: "", Password: ""},
                rules: {
                    Username: [
                        {required: true, min: 3, message: '账号不能为空，长度至少5位', trigger: 'blur'}
                    ],
                    Password: [
                        {required: true, min: 6, message: '密码不能为空，长度至少6位', trigger: 'blur'}
                    ]
                }
            };
        },
        mounted() {
            Cookies.remove(this.conf.cookieTokenName);
            var un = Cookies.get(this.conf.cookieUserName);
            if (un && un != "" && un.length > 1) {
                this.form.Username = un;
            }
        },
        methods: {
            loadAImg() {
                this.form.CaptchaVal = "";
                this.reqUrl = this.prefix + '/api/LoginCaptcha';
                this.loadCaptcha();
            },
            handleSubmit() {
                if (!this.prefix) {
                    this.prefix = "";
                }
                for (var k in this.err) {
                    this.err[k] = "";
                }
                this.$refs.form.validate((valid) => {
                    if (valid) {
                        this.ajax(this.prefix + this.reqUrl, this.form, (r, th) => {
                            if (r.status) {
                                if (r.status.Code == 0) {
                                    th.setUser(r.status);
                                    th.jump(th.goUrl);
                                    return;
                                }
                                if (r.status.Code == 1) {
                                    th.err['Username'] = r.status.msg;
                                }
                                if (r.status.Code == 2) {
                                    th.loadAImg()
                                }
                                if (r.status.Code == 3) {
                                    th.err['Password'] = r.status.msg;
                                    th.form.Password = "";
                                    if (th.authImg != '') {
                                        th.loadAImg();
                                    }
                                }
                            } else {
                                if (r.Code == 8) {
                                    th.form.CaptchaId = r.result[0];
                                    th.form.CaptchaVal = "";
                                    th.authImg = th.prefix + th.authImgUrl + th.form.CaptchaId;
                                }
                            }
                        });
                    }
                });
            },
            loadCaptcha() {
                let th = this;
                this.ajax(this.captchaNew, {}, function (r) {
                    th.CaptchaId = r.result[0];
                    th.authImg = th.prefix + th.authImgUrl + th.CaptchaId;
                });
            }
        }
    };
</script>