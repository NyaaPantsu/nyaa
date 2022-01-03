// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
function Translations() {
  var translations = {}
  var noError = false
  this.Add = function (tr, val) {
    var trans = {}
    if (val != undefined) {
      trans[tr] = val
    } else {
      trans = tr
    }
    Object.assign(translations, trans)
  }
  this.r = function(string, ...args) {
    if ((string != undefined) && (translations[string] != undefined)) {
      if (args != undefined) {
        return this.format(translations[string], ...args)
      }
      return translations[string]
    }
    if (!noError) {
      console.error("No translation string for %s! Please check!", string)
    } else {
      noError = false
    }
    return string
  }
  this.format = function (format, ...args) {
    return format.replace(/{(\d+)}/g, function (match, number) {
      return typeof args[number] != 'undefined' ?
        args[number] :
        match
    })
  }
  this.R = function(string, ...args) {
    noError = true
    return this.r(string, ...args)
  }
}

var T = new Translations()
// @license-end
