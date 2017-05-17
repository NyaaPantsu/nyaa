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