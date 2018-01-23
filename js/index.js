ELEMENT.locale(ELEMENT.lang.ja);
new Vue({
    el: '#main',
    delimiters: ['%%', '%%'],
    data() {
        return {
            activeIndex: '1',
            categories: []
        };
    },
    created: function () {
        var self = this;
        self.getCategories();
    },
    methods: {
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
                        confirmButtonText: 'OK',
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
                    confirmButtonText: 'OK',
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