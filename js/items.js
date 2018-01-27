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
        onFileChange: function (e) {
            var files = e.target.files || e.dataTransfer.files;
            if (!files.length) {
                return;
            }
            this.createImage(files[0]);
        },
        createImage(file) {
            var image = new Image();
            var reader = new FileReader();
            var vm = this;
            reader.onload = function (e) {
                vm.item.item_image = e.target.result;
            };
            reader.readAsDataURL(file);
        },
        removeImage: function (e) {
            this.item.item_image = '';
        }
    }
})