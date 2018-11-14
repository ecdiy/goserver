import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [
        {path: '/', name: 'home', component: () => import('./views/Home.vue')},
        {path: '/rule', name: 'rule', component: () => import('./views/rule.vue')},

        {path: '/user/login', name: 'login', component: () => import('./views/base/login.vue')},
        {path: '/user/register', name: 'login', component: () => import('./views/base/register.vue')},


    ]
})
