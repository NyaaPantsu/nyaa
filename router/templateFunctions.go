package router

import (
	"html/template"
	"log"
	"math"
	"net/url"
	"strconv"
)

var FuncMap = template.FuncMap{
	"min": math.Min,
	"genRoute": func(name string, params ...string) string {
		url, err := Router.Get(name).URL(params...)
		if err == nil {
			return url.String()
		}
		return "error"
	},
	"genNav": func(nav Navigation, currentUrl *url.URL, pagesSelectable int) template.HTML {
		maxPages := math.Ceil(float64(nav.TotalItem) / float64(nav.MaxItemPerPage))

		var ret = ""
		if nav.CurrentPage-1 > 0 {
			url, _ := Router.Get(nav.Route).URL("page", "1")
			ret = ret + "<li><a id=\"page-prev\" href=\"" + url.String() + "?" + currentUrl.RawQuery + "\" aria-label=\"Previous\"><span aria-hidden=\"true\">&laquo;</span></a></li>"
		}
		startValue := 1
		if nav.CurrentPage > pagesSelectable {
			startValue = (int(math.Min((float64(nav.CurrentPage)+math.Floor(float64(pagesSelectable)/2)), maxPages)) - pagesSelectable + 1)
		}
		endValue := (startValue + pagesSelectable - 1)
		if endValue > int(maxPages) {
			endValue = int(maxPages)
		}
		log.Println(nav.TotalItem)
		for i := startValue; i <= endValue; i++ {
			pageNum := strconv.Itoa(i)
			url, _ := Router.Get(nav.Route).URL("page", pageNum)
			ret = ret + "<li"
			if i == nav.CurrentPage {
				ret = ret + " class=\"active\""
			}

			ret = ret + "><a href=\"" + url.String() + "?" + currentUrl.RawQuery + "\">" + strconv.Itoa(i) + "</a></li>"
		}
		if nav.CurrentPage < int(maxPages) {
			url, _ := Router.Get(nav.Route).URL("page", strconv.Itoa(nav.CurrentPage+1))
			ret = ret + "<li><a id=\"page-next\" href=\"" + url.String() + "?" + currentUrl.RawQuery + "\" aria-label=\"Next\"><span aria-hidden=\"true\">&raquo;</span></a></li>"
		}
		return template.HTML(ret)
	},
}
