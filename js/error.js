ELEMENT.locale(ELEMENT.lang.ja);
new Vue({
    el: '#main',
    delimiters: ['%%', '%%'],
    data() {
        return {
            activeIndex: '1',
            redirectTime: 10
        };
    },
    created: function () {
        var self = this;
        self.redirectPage();
    },
    methods: {
        redirectPage: function () {
            var self = this;
            $("#error-description .el-alert__description").text(self.redirectTime + "秒後にリダイレクトします...");
            if (self.redirectTime > 0) {
                self.redirectTime = --self.redirectTime;
            } else {
                location.href = "/login";
            }
            setTimeout(self.redirectPage, 1000)
        }
    }
})