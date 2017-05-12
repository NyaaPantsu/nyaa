var night = localStorage.getItem("night");
function toggleNightMode() {
    var night = localStorage.getItem("night");
    if(night == "true") {
        document.getElementById("style-dark").remove()
    } else {
        document.getElementsByTagName("head")[0].append(darkStyleLink);
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
function formatDate(date) { // thanks stackoverflow
    var monthNames = [
        "January", "February", "March",
        "April", "May", "June", "July",
        "August", "September", "October",
        "November", "December"
    ];

    var day = date.getDate();
    var monthIndex = date.getMonth();
    var year = date.getFullYear();

    return day + ' ' + monthNames[monthIndex] + ' ' + year;
}

var list = document.getElementsByClassName("date-short");
for(var i in list) {
	var e = list[i];
	e.title = e.innerText;
	e.innerText = formatDate(new Date(e.innerText));
}

var list = document.getElementsByClassName("date-full");
for(var i in list) {
	var e = list[i];
	e.title = e.innerText;
	var date = new Date(e.innerText);
	e.innerText = date.toDateString() + " " + date.toLocaleTimeString();
}

/*Fixed-Navbar offset fix*/
window.onload = function() {
  var shiftWindow = function() { scrollBy(0, -70) };
if (location.hash) shiftWindow();
window.addEventListener("hashchange", shiftWindow);
};
