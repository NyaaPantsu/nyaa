var Query = {
    Get: function(url, renderer, callback) {
        var xhr = new XMLHttpRequest();
        console.log(url)
        xhr.open('GET', url, true);
        xhr.responseType = 'json';
        xhr.onload = function(e) {
            if (this.status == 200) {
                renderer(this.response);
                callback(this.response)
            }
        };
        xhr.send();
    }
}