var explosion = document.getElementById("explosion");
var nyanpassu = document.getElementById("nyanpassu");

/*function toggleNightMode() {
	var night = localStorage.getItem("night");
	if(night == "true") {
		document.getElementsByTagName("head")[0].removeChild(darkStyleLink);
	} else {
		document.getElementsByTagName("head")[0].appendChild(darkStyleLink);
	}
	localStorage.setItem("night", (night == "true") ? "false" : "true");
}

function changeTheme(opt) {
	theme = opt.value;
	localStorage.setItem("theme", theme);
	document.getElementById("theme").href = "/css/" + theme;
	console.log(theme);
}

*/


function switchThemes(themeName=null){
	// Switches between themes when a new one is selected

	// Get the theme from the selector if we are not manually passing one
	if(!themeName){
		themeName = document.getElementById("theme-selector").value
	}
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

function toggleNightMode(currentTheme){
	// Turns night mode on and off
	nightMode = sessionStorage.nightMode;
	if(nightMode === undefined){
		nightMode = "false";
	}
	console.log(nightMode);
	// If nightmode is on, turn it off
	if(nightMode === "true"){
		console.log("turning off!")
		sessionStorage.nightMode = "false";
		switchThemes(sessionStorage.previousTheme);
	}
	// If nightmode is turned off, turn it on
	else{
		console.log("turning on!");
		sessionStorage.nightMode = "true";
		// Hardcoded to g for now, I need to find a better solution here
		sessionStorage.setItem("previousTheme", currentTheme);
		switchThemes("tomorrow");
	}
}

function themeFixes(){
	// Loads night mode if it is turned on, meant to be called onload
	// Unloads the mascot too
	if(sessionStorage.nightMode === "true"){
		switchThemes("tomorrow");
	}
	if(localStorage.fuckRenge === "true"){
		document.getElementById("mascot").className = "hide";
	}
}


function disableNightMode(){
	// Called when a user switches themes, we don't want him to automatically go back to nightmode
	sessionStorage.nightMode = "false";
}

function toggleMascot() {
	var fuckRenge = localStorage.fuckRenge;
        if(fuckRenge === undefined){
		fuckRenge = "false";
	}
	// If it is not specifically disabled, covers the undefined case as well
	if (fuckRenge === "false") {
		localStorage.fuckRenge = "true";
		document.getElementById("mascot").className = "hide";
	} else {
		localStorage.fuckRenge = "false";
		document.getElementById("mascot").className = "";
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
window.onload = function() {
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
};

// $(document).ready equivilent, prevents night mode flickering
document.addEventListener("DOMContentLoaded", function(event) { 
	themeFixes();
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