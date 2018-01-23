ELEMENT.locale(ELEMENT.lang.ja);
new Vue({
    el: '#main',
    delimiters: ['%%', '%%'],
    data() {
        return {
            activeIndex: '1',
            category: {
                category_id: null,
                category_name: null
            },
            items: []
        };
    },
    created: function () {
        var self = this;
        self.getItems();
    },
    methods: {
        setSrc: function (img) {
            return "/sample/" + img;
        },
        goEdit: function (id) {
            location.href = "/items?id=" + id;
        },
        getItems: function () {
            var self = this;
            var item = {
                "item_id": "1",
                "item_name": "A",
                "item_img": "omen_hannya.png"
            }
            self.items.push(item);
            self.items.push(item);
        }
    }
})