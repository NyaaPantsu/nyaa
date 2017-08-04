// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
function Translations() {
  var translations = {};
  this.Add = function (tr, val) {
    trans = {}
    if (val != undefined) {
      trans[tr] = val;
    } else {
      trans = tr
    }
    Object.assign(translations, trans);
  };
  this.r = function (string, ...args) {
    if ((string != undefined) && (translations[string] != undefined)) {
      if (args != undefined) {
        return this.format(translations[string], ...args);
      }
      return translations[string];
    }
    console.error("No translation string for %s! Please check!", string);
    return "";
  };
  this.format = function (format, ...args) {
    return format.replace(/{(\d+)}/g, function (match, number) {
      return typeof args[number] != 'undefined' ?
        args[number] :
        match;
    });
  };
}

var T = new Translations();
// @license-end
