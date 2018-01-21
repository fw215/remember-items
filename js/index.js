
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
            var category = {
                "category_id": "1",
                "category_name": "A",
            }
            self.categories.push(category);
            self.categories.push(category);
        }
    }
})