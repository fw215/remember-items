ELEMENT.locale(ELEMENT.lang.ja);
new Vue({
    el: '#main',
    delimiters: ['%%', '%%'],
    data() {
        return {
            activeIndex: '1',
            categoryModalVisible: false,
            formLabelWidth: '85px',
            categories: [],
            category: {
                'category_id': 0,
                'category_name': '',
            }
        };
    },
    created: function () {
        var self = this;
        self.getCategories();
    },
    methods: {
        addCategory: function () {
            var self = this;
            var params = new URLSearchParams();
            params.append('category_id', self.category.category_id);
            params.append('category_name', self.category.category_name);
            axios.post(
                "/v1/category",
                params
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.category = {
                        'category_id': 0,
                        'category_name': '',
                    };
                    self.getCategories();
                    self.categoryModalVisible = false;
                } else {
                    self.$message({
                        dangerouslyUseHTMLString: true,
                        message: res.data.errors.join(`<br>`),
                        type: 'error'
                    });
                }
            }).catch(function (error) {
                self.categoryModalVisible = false;
                self.$message({
                    type: 'warning',
                    message: `しばらく時間をおいてから再度お試しください`
                });
            });
        },
        goEdit: function (id) {
            location.href = "/items?id=" + id;
        },
        getCategories: function () {
            var self = this;
            self.error_message = "";
            axios.get(
                "/v1/categories"
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.categories = res.data.categories;
                } else {
                    self.$alert('エラーが発生しました', 'エラー', {
                        confirmButtonText: 'はい',
                        callback: function (action) {
                            self.$message({
                                type: 'warning',
                                message: res.data.error
                            });
                        }
                    });
                }
            }).catch(function (error) {
                self.$alert('通信エラーが発生しました', 'エラー', {
                    confirmButtonText: 'はい',
                    callback: function (action) {
                        self.$message({
                            type: 'warning',
                            message: `しばらく時間をおいてから再度お試しください`
                        });
                    }
                });
            });
        }
    }
})