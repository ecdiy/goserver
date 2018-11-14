import Vue from 'vue'
import App from './App.vue'
import router from './router'
import iView from 'iview';
import Cookies from 'js-cookie';
import 'iview/dist/styles/iview.css';
import web from '@/components/web/index';

Vue.use(iView);
Vue.use(web, {
    cookieTokenName: 'adminToken', cookieUserName: 'adminUser',
    cookieUserInfoName: 'adminInfo'
});
Vue.config.productionTip = false

Vue.prototype.setUser = function (obj) {
    Cookies.set(this.conf.cookieTokenName, obj.Token, {expires: 365});
    delete obj.Token;
    delete obj.code;
    Cookies.set(this.conf.cookieUserName, obj.Username, {expires: 365});
    Cookies.set(this.conf.cookieUserInfoName, JSON.stringify(obj), {expires: 365});
}
window.vm = new Vue({
    router,
    render: h => h(App)
}).$mount('#app')
