<template>
    <div>
        <list url="/gk-admin/sp/UserList" :columns="columns"/>
    </div>
</template>

<script>
    import list from '@/components/web/list'

    export default {
        data(){
            return {columns:[
                    {title: 'Id', key: 'Id'},
                    {title: 'Username', key: 'Username'},
                ]}
        },
        components: {list}
    };
</script>

<style scoped>

</style>