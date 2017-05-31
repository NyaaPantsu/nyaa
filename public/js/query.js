var Query = {
    Failed:0,
    MaxFail: -1,
    Get: function(url, renderer, callback) {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', url, true);
        xhr.responseType = 'json';
        xhr.onload = function(e) {
            if (this.status == 200) {
                Query.Failed = 0;
                renderer(this.response);
                if (callback != undefined) callback(this.response);
            } else {
                console.log("Error when refresh")
                Query.Failed++;
                console.log("Attempt to refresh "+Query.Failed+"...");
                if ((Query.MaxFail == -1) || (Query.Failed < Query.MaxFail)) Query.Get(url, renderer, callback);
                else console.error("Too many attempts, stopping...")
            }
        };
        xhr.send();
    }
};