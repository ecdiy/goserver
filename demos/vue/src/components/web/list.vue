<template>
    <card>
        <Table :columns="columns" :loading="loading" :data="list"></Table>
        <div style="clear: both;padding-top: 10px; ">
            <slot style="float: left"></slot>
            <Page   :total="total" style="text-align: right; padding-top: 10px; clear: both;"
                  :current="page" :page-size="pageSize" show-sizer show-total @on-change="load"></Page>
        </div>
    </card>
</template>
<script>
    export default {
        props: ['columns', 'url', 'param'],
        data() {
            return {total: 0, pageSize: 20, page: 1, list: [], loading: true}
        },

        mounted() {
            this.page = Number(this.$route.params.page ? this.$route.params.page : '1');
            this.loadX();
        },

        methods: {
            load(data) {
                this.page = data;
                this.loadX()
            },
            loadX() {
                var p;
                if (this.param) {
                    p = this.param;
                } else {
                    p = {}
                }
                p.Begin = (this.page - 1) * this.pageSize;
                this.ajax(this.url, p)
            }
        }
    }
</script>