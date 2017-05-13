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