// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
document.querySelector("[data-selectall='checkbox']").addEventListener("change", function(e) {
  var cbs = document.querySelectorAll("input[type='checkbox'].selectable");
  var l = cbs.length;
  for (var i=0; i<l; i++) cbs[i].checked = e.target.checked;
});
// @license-end
