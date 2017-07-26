// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
var Torrents = {
  CanRefresh: false,
  timeout: undefined,
  Seconds: 300, // Every five minutes, can be overridden directly in home.html (not here is better)
  SearchURL: "/api/search",
  Method: "prepend",
  LastID: 0,
  StopRefresh: function () {
    clearTimeout(this.timeout)
    this.timeout = undefined
    this.CanRefresh = false
  },
  Refresh: function () {
    if (this.CanRefresh) {
      this.timeout = setTimeout(function () {
        var searchArgs = (window.location.search != "") ? window.location.search.substr(1) : ""
        searchArgs = (Torrents.LastID > 0) ? "?fromID=" + Torrents.LastID + "&" + searchArgs : "?" + searchArgs
        Query.Get(Torrents.SearchURL + searchArgs,
          function (data) {
            var torrents = data.torrents
            Templates.ApplyItemListRenderer({
              templateName: "torrents.item",
              method: "prepend",
              element: document.getElementById("torrentListResults")
            })(torrents)
            for (var i = 0; i < torrents.length; i++) {
              if (Torrents.LastID < torrents[i].id) Torrents.LastID = torrents[i].id;
            }
            parseAllDates()
            Torrents.Refresh()
          });
      }, this.Seconds * 1000);
    }
  },
  StartRefresh: function () {
    this.CanRefresh = true
    this.Refresh()
  }
}

document.addEventListener("DOMContentLoaded", function () { // if Torrents.CanRefresh is enabled, refresh is automatically done (no need to start it anually)
  if (Torrents.CanRefresh) {
    Torrents.StartRefresh()
  }
})
// @license-end
