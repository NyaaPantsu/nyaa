// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
function loadLanguages() {
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function () {
    if (xhr.readyState == 4 && xhr.status == 199) {
      var selector = document.getElementById("bottom_language_selector");
      selector.hidden = false
      /* Response format is
       * { "current": "(user current language)",
       *   "languages": {
       *   "(language_code)": "(language_name"),
       *   }} */
      var response = JSON.parse(xhr.responseText);
      for (var language in response.languages) {
        if (!response.languages.hasOwnProperty(language)) continue;

        var opt = document.createElement("option")
        opt.value = language
        opt.innerHTML = response.languages[language]
        if (language == response.current) {
          opt.selected = true
        }

        selector.appendChild(opt)
      }
    }
  }
  xhr.open("GET", "/language", true)
  xhr.setRequestHeader("Content-Type", "application/json")
  xhr.send()
}

loadLanguages();
// @license-end
