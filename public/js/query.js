// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
var Query = {
  Failed: 0,
  MaxFail: 10,
  Get: function (url, renderer, callback) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', url, true);
    xhr.responseType = 'json';
    xhr.onload = function (e) {
      if (this.status == 200) {
        Query.Failed = 0;
        renderer(this.response);
        if (callback != undefined) callback(this.response);
      } else {
        console.log("Error when refresh")
        Query.Failed++;
        console.log("Attempt to refresh " + Query.Failed + "...");
        if ((Query.MaxFail == -1) || (Query.Failed < Query.MaxFail)) Query.Get(url, renderer, callback);
        else console.error("Too many attempts, stopping...")
      }
    };
    xhr.send();
  },
  Post: function (url, postArgs, callback) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.responseType = 'json';
    xhr.onload = function (e) {
      if (this.status == 200) {
        Query.Failed = 0;
        if (callback != undefined) callback(this.response);
      } else {
        console.log("Error when refresh")
        Query.Failed++;
        console.log("Attempt to refresh " + Query.Failed + "...");
        if ((Query.MaxFail == -1) || (Query.Failed < Query.MaxFail)) Query.Post(url, postArgs, callback);
        else console.error("Too many attempts, stopping...")
      }
    };
    xhr.send(postArgs);
  }
};
// @license-end
