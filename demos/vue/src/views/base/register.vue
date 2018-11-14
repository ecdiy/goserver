<template>
    <div>
        <card>
            <Form :label-width="80" style="width:400px;margin: auto">
                <FormItem prop="Username" label="用户名">
                    <Input v-model="reg.Username" placeholder="请输入用户名">
                    <span slot="prepend"><Icon :size="16" type="person"></Icon></span></Input>
                    <div class="ivu-form-item-error-tip" v-show="em.Username!=''">{{em.Username}}</div>
                </FormItem>
                <FormItem prop="Password" label="密码">
                    <Input v-model="reg.Password" placeholder="请输入密码">
                    <span slot="prepend"><Icon :size="16" type="person"></Icon></span></Input>
                    <div class="ivu-form-item-error-tip" v-show="em.Password!=''">{{em.Password}}</div>
                </FormItem>
                <FormItem prop="CaptchaVal" label="认证码">
                    <Input type="text" v-model="reg.CaptchaVal" placeholder="请输入认证码">
                    <span slot="prepend"><Icon :size="14" type="flash"></Icon></span></Input>
                    <img :src="authImg" @click="loadCaptcha"/>
                    <div class="ivu-form-item-error-tip" v-show="em.CaptchaVal!=''">{{em.CaptchaVal}}</div>
                </FormItem>
                <FormItem>
                    <Button type="primary" long size="large" @click="registerSubmit">注册</Button>
                    <div style="float: right">
                        <router-link to="/user/login">登录</router-link>
                        <router-link to="/user/forget">忘记密码</router-link>
                    </div>
                </FormItem>
            </Form>
        </card>

    </div>
</template>

<script>
    export default {
        data() {
            return {
                prefix: '/gk-admin',
                authImg: "", reg: {},
                em: {Username: '', Email: '', Mobile: '', Password: '', CaptchaVal: ''},
            }
        }, mounted() {
            this.loadCaptcha()
        }, methods: {
            validate() {
                var un = this.reg.Username;
                if (!un || un.length > 32 || un.length < 5) {
                    this.em.Username = '用户名长度5~32';
                    return false;
                } else {
                    this.em.Username = '';
                }
                var e   = this.reg.CaptchaVal;
                if (!e || e.length != 6) {
                    this.em.CaptchaVal = '认证码长度6';
                    return false;
                } else {
                    this.em.CaptchaVal = '';
                }
                return true;
            },
            loadCaptcha() {
                this.ajax(this.prefix + '/api/CaptchaNew', {}, (r, th) => {
                    th.CaptchaId = r.result[0];
                    th.authImg = th.prefix + "/Captcha?t=" + th.CaptchaId;
                });
            },
            registerSubmit() {
                if (!this.validate()) {
                    return
                }
                this.ajax(this.prefix + '/api/Register', this.reg, function (r, th) {
                    if (r.Code == 8) {
                        th.reg.CaptchaId = r.result[0];
                        th.reg.CaptchaVal = "";
                        th.authImg = th.prefix + "/Captcha?t=" + th.reg.CaptchaId;
                        th.em.CaptchaVal = r.msg;
                        return
                    }
                    th.loadCaptcha();
                    if (r.result.Code == 1000) {
                        th.reg.CaptchaVal = "";
                        th.em.Username = r.result.msg;
                        th.loadCaptcha()
                    }

                    if (r.result.Code == 0) {
                        th.$Modal.success({
                            title: "", content: "注册成功",
                            onOk: function () {
                                //  location.href = "/";
                            }
                        });
                    }
                });

            },
        }
    }
</script>

<style scoped>

</style>