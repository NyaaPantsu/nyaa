// Night mode
// also sorry that this code is soo bad, literally nothing else worked.. ima remake it in a near future
var night = localStorage.getItem("night");
if (night=="true") {
    document.getElementById("style").href = "/css/style-night.css";
    document.getElementById("nightbutton").innerHTML = "<img id='sunmoon' src='/img/sun.png' alt='Day!'>";
}

function toggleNightMode() {
    var styleshieeet = document.getElementById("style").href;
    var styleshieet = new RegExp("style.css");
    var stylesheet = styleshieet.test(styleshieeet);
    if (stylesheet==true) {
	document.getElementById("style").href = "/css/style-night.css";
	document.getElementById("nightbutton").innerHTML = "<img id='sunmoon' src='/img/sun.png' alt='Day!'>";
	localStorage.setItem("night", "true");
    }
    else {
	document.getElementById("style").href = "/css/style.css";
	document.getElementById("nightbutton").innerHTML = "<img id='sunmoon' src='/img/moon.png' alt='Night!'>";
	localStorage.setItem("night", "false");
    }
console.log(styleshieeet);
console.log(stylesheet);
}

// Used by spoiler tags
function toggleLayer(elem) {
	if (elem.classList.contains("hide"))
		elem.classList.remove("hide");
	else
		elem.classList.add("hide");
}

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
