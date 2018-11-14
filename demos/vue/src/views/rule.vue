<template>
    <div style="width: 50%;margin: auto">
        <card>
            <Button @click="callSysTest">
                call sp TestSys: /gk-admin/sys/Test : {{testCode}}
            </Button>
            认证失败，没有权限时
        </card>
        <card>
            <div>
                <Button @click="callOk">登录后再调用</Button>
                <div>{{rule}}</div>
            </div>
        </card>

    </div>
</template>

<script>import Cookies from 'js-cookie';
import axios from 'axios';

export default {
    data() {
        return {testCode: 0, username: '', rule: {}}
    },
    methods: {

        callSysTest() {
            let th = this;
            axios.post('/gk-admin/sys/Test', {}).then(function (r) {
            }).catch((err) => {
                console.log("~~~", err.response.status, "~~", err)
                th.testCode = err.response.status;
            })
        },
        callOk() {
            let th = this;
            axios.post('/gk-admin/sp/UserLogin', {Username: 'test', Password: 'test'}).then(function (r) {

                console.log(r.data.status.Token)

                axios.post('/gk-admin/sys/Test', {token: r.data.status.Token}).then(function (r) {
                    console.log(r)
                    th.rule = r.data
                }).catch((err) => {
                    console.log("~~~", err.response.status, "~~", err)
                    th.testCode = err.response.status;
                })


            }).catch((err) => {
                console.log("~~~", err.response.status, "~~", err)
                th.testCode = err.response.status;
            })
        }
    }
}
</script>

<style scoped>

</style>