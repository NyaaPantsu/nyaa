function loadLanguages() {
	var xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function() {
		if (xhr.readyState == 4 && xhr.status == 200) {
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
