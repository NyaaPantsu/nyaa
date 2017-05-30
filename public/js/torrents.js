var Torrents = {
    CanRefresh: false,
    timeout: undefined,
    Seconds: 300, // Every five minutes, can be overridden directly in home.html (not here is better)
    SearchURL: "/api/search",
    Method: "prepend",
    LastID: 0,
    StopRefresh: function() {
        clearTimeout(this.timeout)
        this.timeout = undefined
        this.CanRefresh = false
    },
    Refresh: function() {
        if (this.CanRefresh) {
            console.log("Start Refresh...")
            this.timeout = setTimeout(function() {
                var searchArgs = (window.location.search != "") ? window.location.search.substr(1) : ""
                searchArgs = (Torrents.LastID > 0) ? "?fromID="+Torrents.LastID+"&"+searchArgs : "?"+searchArgs
                Query.Get(Torrents.SearchURL+searchArgs, 
                    Templates.ApplyItemListRenderer({
                        templateName: "torrents.item", method: "prepend", element: document.getElementById("torrentListResults")
                    }), function(torrents) {
                        for (var i =0; i < torrents.length; i++) { if (Torrents.LastID < torrents[i].id) Torrents.LastID = torrents[i].id; }
                        parseAllDates();
                        Torrents.Refresh()
                    });
            }, this.Seconds*1000);
        }
    },
    StartRefresh: function() {
        this.CanRefresh = true;
        this.Refresh()
    }
}

document.addEventListener("DOMContentLoaded", function() { // if Torrents.CanRefresh is enabled, refresh is automatically done (no need to start it anually)
    if (Torrents.CanRefresh) {
        Torrents.StartRefresh()
    }
})


// Credits to mpen (StackOverflow)
function humanFileSize(bytes, si) {
    var thresh = si ? 1000 : 1024;
    if(Math.abs(bytes) < thresh) {
        return bytes + ' B';
    }
    var units = si
        ? ['kB','MB','GB','TB','PB','EB','ZB','YB']
        : ['KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB'];
    var u = -1;
    do {
        bytes /= thresh;
        ++u;
    } while(Math.abs(bytes) >= thresh && u < units.length - 1);
    return bytes.toFixed(1)+' '+units[u];
}