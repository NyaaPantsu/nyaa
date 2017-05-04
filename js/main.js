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
            el.href = pageString + n + query;
            el.innerHTML = n;
        }

        document.getElementById('page-next').href = pageString + next + query;
        document.getElementById('page-prev').href = pageString + prev + query;
