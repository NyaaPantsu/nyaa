var Query = {
    Failed:0,
    MaxConsecutingFailing:-1,
    Get: function(url, renderer, callback) {
        var xhr = new XMLHttpRequest();
        console.log(url)
        xhr.open('GET', url, true);
        xhr.responseType = 'json';
        xhr.onload = function(e) {
            if (this.status == 200) {
                Query.Failed = 0;
                renderer(this.response);
                callback(this.response);
            } else {
                console.log("Error when refresh")
                Query.Failed++;
                console.log("Attempt to refresh "+Query.Failed+"...");
                if ((Query.MaxConsecutingFailing == -1) || (Query.Failed < Query.MaxConsecutingFailing)) Query.Get();
                else console.error("Too many attempts, stopping...")
            }
        };
        xhr.send();
    }
}