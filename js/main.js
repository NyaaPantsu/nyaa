        var pathArray = window.location.pathname.split( '/' );
        var query = window.location.search;
        var page = parseInt(pathArray[2]);
        var pageString = "/page/";

        var next = page + 1;
        var prev = page - 1;

        if (prev < 1) {
            prev = 1;
        }

        if (isNaN(page)) {
            next =  2;
            prev =  1;
        }

        if (query != "") {
            pageString = "/search/";
        }

        var maxId = 5;
        for (var i = 0; i < maxId; i++) {
            var el = document.getElementById('page-' + i), n = prev + i;
            if (el == null)
                continue;
            el.href = pageString + n + query;
            el.innerHTML = n;
        }

        var e = document.getElementById('page-next');
        if (e != null)
            e.href = pageString + next + query;
        var e = document.getElementById('page-prev');
        if (e != null)
            e.href = pageString + prev + query;

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
