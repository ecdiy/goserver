<template>
    <card>
        <div style="width:300px;margin: auto">
            <div class="login-con">
                <Card :bordered="false">

                    <div class="form-con">
                        <Form :model="form" ref="form">

                            <FormItem prop="Password" :rules="rules.Password" label="密码"
                                      :error="err.Password">
                                <Input type="password" v-model="form.Password" placeholder="请输入密码">
                                <span slot="prepend">
                                    <Icon :size="14" type="locked"></Icon>
                                </span>
                                </Input>
                            </FormItem>

                            <FormItem>
                                <Button @click="handleSubmit" type="primary" long>修改</Button>
                            </FormItem>
                        </Form>
                    </div>
                </Card>
            </div>
        </div>
    </card>
</template>


<script>

    export default {
        data() {
            return {

                form: {},
                err: {Password: ""},
                rules: {

                    Password: [
                        {required: true, min: 6, message: '密码不能为空，长度至少6位', trigger: 'blur'}
                    ]
                }
            };
        },

        methods: {
            handleSubmit() {
                for (var k in this.err) {
                    this.err[k] = "";
                }
                this.$refs.form.validate((valid) => {
                    if (valid) {
                        this.ajax('/gk-admin/sp/UserUpPassword', this.form, (r, th) => {
                            th.$Notice.success({title: '设置保存成功'});
                        });
                    }
                });
            }
        }
    };
</script>