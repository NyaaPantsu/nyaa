var night = localStorage.getItem("night");
function toggleNightMode() {
    var night = localStorage.getItem("night");
    if(night == "true") {
        document.getElementsByTagName("head")[0].removeChild(darkStyleLink);
    } else {
        document.getElementsByTagName("head")[0].appendChild(darkStyleLink);
    }
    localStorage.setItem("night", (night == "true") ? "false" : "true");
}

// Used by spoiler tags
function toggleLayer(elem) {
	if (elem.classList.contains("hide"))
		elem.classList.remove("hide");
	else
		elem.classList.add("hide");
}

// Date formatting
var lang = $("html").attr("lang");
var shortOpt = { year: "numeric", month: "short", day: "numeric" };

var list = document.getElementsByClassName("date-short");
for(var i in list) {
	var e = list[i];
	e.title = e.innerText;
	e.innerText = new Date(e.innerText).toLocaleString(lang, shortOpt);
}

var list = document.getElementsByClassName("date-full");
for(var i in list) {
	var e = list[i];
	e.title = e.innerText;
	e.innerText = new Date(e.innerText).toLocaleString(lang);
}

/*Fixed-Navbar offset fix*/
window.onload = function() {
  var shiftWindow = function() { scrollBy(0, -70) };
if (location.hash) shiftWindow();
window.addEventListener("hashchange", shiftWindow);
};
function loadLanguages() {
	var xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function() {
		if (xhr.readyState == 4 && xhr.status == 200) {
			var selector = document.getElementById("bottom_language_selector");
			selector.hidden = false
			/* Response format is
			 * { "current": "(user current language)",
			 *   "languages": {
			 *   	"(language_code)": "(language_name"),
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
	xhr.open("GET", "/language?format=json", true)
	xhr.send()
}

loadLanguages();
