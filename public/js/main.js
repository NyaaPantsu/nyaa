var explosion = document.getElementById("explosion");
var nyanpassu = document.getElementById("nyanpassu");

function toggleNightMode() {
	var night = localStorage.getItem("night");
	if(night == "true") {
		document.getElementsByTagName("head")[0].removeChild(darkStyleLink);
	} else {
		document.getElementsByTagName("head")[0].appendChild(darkStyleLink);
	}
	localStorage.setItem("night", (night == "true") ? "false" : "true");
}

// Switches between themes when a new one is selected
function switchThemes(){
	themeName = document.getElementById("theme-selector").value
	var head = document.getElementsByTagName("head")[0];
	// Remove the theme in place, it fails if one isn't set
	try{
		head.removeChild(document.getElementById("theme"));
	} catch(err){}
	// Don't add a node if we don't want extra styling
	if(themeName === ""){
		return;
	}
	// Create the new one and put it back
        var newTheme = document.createElement("link");
        newTheme.setAttribute("rel", "stylesheet");
        newTheme.setAttribute("href", "/css/"+ themeName + ".css");
        newTheme.setAttribute("id", "theme");
	head.appendChild(newTheme);
}


function toggleMascot(btn) {
	var state= btn.value;
	if (state == "hide") {
		btn.innerHTML = "Mascot";
		document.getElementById("mascot").className = "hide";
		btn.value = "show";
	} else {
		btn.innerHTML = "Mascot";
		document.getElementById("mascot").className = "";
		btn.value = "hide";
	}
}

// Used by spoiler tags
function toggleLayer(elem) {
	if (elem.classList.contains("hide"))
		elem.classList.remove("hide");
	else
		elem.classList.add("hide");
}

// Date formatting
var lang = document.getElementsByTagName("html")[0].getAttribute("lang"); 
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
/*Fixed-Navbar offset fix*/
document.addEventListener("DOMContentLoaded", function(event) {
	var shiftWindow = function() { scrollBy(0, -70) };
	if (location.hash) shiftWindow();
	window.addEventListener("hashchange", shiftWindow);
	
	document.getElementsByClassName("search-box")[0].addEventListener("focus", function (e) {
		var w = document.getElementsByClassName("h-user")[0].offsetWidth;
		document.getElementsByClassName("h-user")[0].style.display = "none";
		document.getElementsByClassName("search-box")[0].style.width = document.getElementsByClassName("search-box")[0].offsetWidth + w + "px";
	});
	document.getElementsByClassName("search-box")[0].addEventListener("blur", function (e) {
		document.getElementsByClassName("search-box")[0].style.width = "";
		document.getElementsByClassName("h-user")[0].style.display = "inline-block";
	});

	// Keep mascot hiding choice
	var hideMascot = (localStorage.getItem("hide_mascot") == "true")
	if (hideMascot) {
		var btn = document.getElementById("mascotKeepHide");
		btn.innerHTML = "Mascot";
		document.getElementById("mascot").className = "hide";
		btn.value = "show";
	}
});

function playVoice() {
	if (explosion) {
		explosion.play();
	}
	else {
		nyanpassu.volume = 0.5;
		nyanpassu.play();
	}
}

function toggleMascot(btn) {
	var state= btn.value;
	if (state == "hide") {
		btn.innerHTML = "Mascot";
		document.getElementById("mascot").className = "hide";
		btn.value = "show";
		localStorage.setItem("hide_mascot", "true")
	} else {
		btn.innerHTML = "Mascot";
		document.getElementById("mascot").className = "";
		btn.value = "hide";
		localStorage.setItem("hide_mascot", "false")
	}
}