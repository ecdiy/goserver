import Cookies from 'js-cookie';
import axios from 'axios';

const conf = {
    cookieTokenName: "token", cookieUserName: 'user', cookieUserInfoName: 'userInfo',
    loginUrl: '/user/login'
}

export default {

    install(Vue, options) {
        if (options) {
            for (var k in options) {
                conf[k] = options[k]
            }
        }
        Vue.mixin({
            computed: {
                login() {
                    var tk = Cookies.get(conf.cookieTokenName)
                    if (tk && tk.length > 5) {
                        return true
                    }
                    return false
                },
                user() {
                    var us = Cookies.get(conf.cookieUserInfoName)
                    if (us && us.length > 2) {
                        try {
                            return JSON.parse(us)
                        } catch (e) {
                        }
                    }
                    return {}
                }
            }
        })

        Vue.prototype.conf = conf

        Vue.prototype.jump = function (url) {
            location.href = url
        }

        Vue.prototype.ajax = function (url, p, fun) {
            let th = this;
            axios.post(url, p ? p : {}).then(function (r) {
                for (var k in r.data) {
                    if (th.hasOwnProperty(k))
                        th[k] = r.data[k]
                }
                if (th.hasOwnProperty("loading")) {
                    th.loading = false;
                }
                if (fun) {
                    if (typeof(fun) == 'function') {
                        fun(r.data, th)

                    }
                    if (typeof(fun) == 'string') {
                        th.$router.replace('/' + goUrl)
                    }
                }
            }).catch((err) => {
                if (err.response && err.response.status == 401) {
                    Cookies.remove(conf.cookieTokenName);
                    th.$router.replace(conf.loginUrl)
                } else {
                    console.log(err)
                }
            })
        }
    }
};