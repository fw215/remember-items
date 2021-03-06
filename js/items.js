ELEMENT.locale(ELEMENT.lang.ja); Vue.config.devtools = true;
new Vue({
    el: '#main',
    delimiters: ['%%', '%%'],
    data() {
        return {
            activeIndex: '1',
            itemModalVisible: false,
            formLabelWidth: '70px',
            category: {
                category_id: null,
                category_name: null
            },
            items: [],
            item: {
                item_id: 0,
                item_name: null,
                item_image: null
            }
        };
    },
    created: function () {
        var self = this;
        self.category.category_id = document.getElementById('main').getAttribute('data-category-id');
        self.getCategory();
        self.getItems();
    },
    methods: {
        registerCategory: function () {
            var self = this;
            var params = new URLSearchParams();
            params.append('category_id', self.category.category_id);
            params.append('category_name', self.category.category_name);
            axios.post(
                "/v1/category",
                params
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.$message({
                        message: `カテゴリ名を更新しました`,
                        type: 'success'
                    });
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
        addOpenItem: function () {
            var self = this;
            self.item = {
                item_id: 0,
                item_name: null,
                item_image: null
            };
            self.itemModalVisible = true;
        },
        openItem: function (id) {
            var self = this;
            axios.get(
                "/v1/item/" + id
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.item = res.data.item;
                    self.itemModalVisible = true;
                } else {
                    self.$message({
                        dangerouslyUseHTMLString: true,
                        message: res.data.errors.join(`<br>`),
                        type: 'error'
                    });
                }
            }).catch(function (error) {
                self.itemModalVisible = false;
                self.$message({
                    type: 'warning',
                    message: `しばらく時間をおいてから再度お試しください`
                });
            });
        },
        deleteItem: function () {
            var self = this;
            axios.delete(
                "/v1/item/" + self.item.item_id,
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.item = {
                        item_id: 0,
                        item_name: null,
                        item_image: null
                    };
                    self.$message({
                        message: `アイテムを削除しました`,
                        type: 'success'
                    });
                    self.getItems();
                    self.itemModalVisible = false;
                } else {
                    self.$message({
                        dangerouslyUseHTMLString: true,
                        message: res.data.errors.join(`<br>`),
                        type: 'error'
                    });
                }
            }).catch(function (error) {
                self.itemModalVisible = false;
                self.$message({
                    type: 'warning',
                    message: `しばらく時間をおいてから再度お試しください`
                });
            });
        },
        registerItem: function () {
            var self = this;
            var params = new URLSearchParams();
            params.append('category_id', self.category.category_id);
            params.append('item_id', self.item.item_id);
            params.append('item_name', self.item.item_name);
            params.append('item_image', self.item.item_image);
            axios.post(
                "/v1/item",
                params
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.item = {
                        item_id: 0,
                        item_name: null,
                        item_image: null
                    };
                    self.getItems();
                    self.itemModalVisible = false;
                } else {
                    self.$message({
                        dangerouslyUseHTMLString: true,
                        message: res.data.errors.join(`<br>`),
                        type: 'error'
                    });
                }
            }).catch(function (error) {
                self.itemModalVisible = false;
                self.$message({
                    type: 'warning',
                    message: `しばらく時間をおいてから再度お試しください`
                });
            });
        },
        setSrc: function (img) {
            return "/images/" + img;
        },
        getCategory: function () {
            var self = this;
            self.error_message = "";
            axios.get(
                "/v1/category/" + self.category.category_id
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.category = res.data.category;
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
        },
        getItems: function () {
            var self = this;
            self.error_message = "";
            axios.get(
                "/v1/items/" + self.category.category_id
            ).then(function (res) {
                if (res.data.code == 200) {
                    self.items = res.data.items;
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
        },
        goIndex: function () {
            location.href = "/";
        },
    }
})