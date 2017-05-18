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
var ymdOpt = { year: "numeric", month: "short", day: "numeric" };
var hmOpt  = { hour: "numeric", minute: "numeric" };

var list = document.getElementsByClassName("date-short");
for(var i in list) {
	var e = list[i];
	e.title = e.innerText;
	e.innerText = new Date(e.innerText).toLocaleString(lang, ymdOpt);
}

var list = document.getElementsByClassName("date-full");
for(var i in list) {
	var e = list[i];
	e.title = e.innerText;
	e.innerText = new Date(e.innerText).toLocaleString(lang);
}
$(".date-comments").each(function(index, el) {
	$(this).attr("title", el.innerText);
	$(this).text(new Date($(this).attr("title")).toLocaleDateString(lang, ymdOpt) + " ");
	$(this).append($('<span class="hidden-xs"></span>').text(new Date($(this).attr("title")).toLocaleTimeString(lang, hmOpt)))
});
/*Fixed-Navbar offset fix*/
window.onload = function() {
  var shiftWindow = function() { scrollBy(0, -70) };
if (location.hash) shiftWindow();
window.addEventListener("hashchange", shiftWindow);
};
