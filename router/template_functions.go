package router

import (
	"html/template"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/categories"
	"github.com/NyaaPantsu/nyaa/util/filelist"
	"github.com/NyaaPantsu/nyaa/util/languages"
)

type captchaData struct {
	CaptchaID string
	T         languages.TemplateTfunc
}

// FuncMap : Functions accessible in templates by {{ $.Function }}
var FuncMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"min": math.Min,
	"genRoute": func(name string, params ...string) string {
		url, err := Router.Get(name).URL(params...)
		if err == nil {
			return url.String()
		}
		return "error"
	},
	"genRouteWithQuery": func(name string, currentUrl *url.URL, params ...string) template.URL {
		url, err := Router.Get(name).URL(params...)
		if err == nil {
			return template.URL(url.String() + "?" + currentUrl.RawQuery)
		}
		return "error"
	},
	"genViewTorrentRoute": func(torrent_id uint) string {
		// Helper for when you have an uint while genRoute("view_torrent", ...) takes a string
		// FIXME better solution?
		s := strconv.FormatUint(uint64(torrent_id), 10)
		url, err := Router.Get("view_torrent").URL("id", s)
		if err == nil {
			return url.String()
		}
		return "error"
	},
	"genSearchWithOrdering": func(currentUrl url.URL, sortBy string) template.URL {
		values := currentUrl.Query()
		order := false //Default is DESC
		sort := "2"    //Default is Date (Actually ID, but Date is the same thing)

		if _, ok := values["order"]; ok {
			order, _ = strconv.ParseBool(values["order"][0])
		}
		if _, ok := values["sort"]; ok {
			sort = values["sort"][0]
		}

		if sort == sortBy {
			order = !order //Flip order by repeat-clicking
		} else {
			order = false //Default to descending when sorting by something new
		}

		values.Set("sort", sortBy)
		values.Set("order", strconv.FormatBool(order))

		url, _ := Router.Get("search").URL()
		url.RawQuery = values.Encode()

		return template.URL(url.String())
	},
	"genSortArrows": func(currentUrl url.URL, sortBy string) template.HTML {
		values := currentUrl.Query()
		leftclass := "sortarrowdim"
		rightclass := "sortarrowdim"

		order := false
		sort := "2"

		if _, ok := values["order"]; ok {
			order, _ = strconv.ParseBool(values["order"][0])
		}
		if _, ok := values["sort"]; ok {
			sort = values["sort"][0]
		}

		if sort == sortBy {
			if order {
				rightclass = ""
			} else {
				leftclass = ""
			}
		}

		arrows := "<span class=\"sortarrowleft " + leftclass + "\">▼</span><span class=\"" + rightclass + "\">▲</span>"

		return template.HTML(arrows)
	},
	"genNav": func(nav navigation, currentUrl *url.URL, pagesSelectable int) template.HTML {
		var ret = ""
		if nav.TotalItem > 0 {
			maxPages := math.Ceil(float64(nav.TotalItem) / float64(nav.MaxItemPerPage))

			if nav.CurrentPage-1 > 0 {
				url, _ := Router.Get(nav.Route).URL("page", "1")
				ret = ret + "<li><a id=\"page-prev\" href=\"" + url.String() + "?" + currentUrl.RawQuery + "\" aria-label=\"Previous\"><span aria-hidden=\"true\">&laquo;</span></a></li>"
			}
			startValue := 1
			if nav.CurrentPage > pagesSelectable/2 {
				startValue = (int(math.Min((float64(nav.CurrentPage)+math.Floor(float64(pagesSelectable)/2)), maxPages)) - pagesSelectable + 1)
			}
			if startValue < 1 {
				startValue = 1
			}
			endValue := (startValue + pagesSelectable - 1)
			if endValue > int(maxPages) {
				endValue = int(maxPages)
			}
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
			itemsThisPageStart := nav.MaxItemPerPage*(nav.CurrentPage-1) + 1
			itemsThisPageEnd := nav.MaxItemPerPage * nav.CurrentPage
			if nav.TotalItem < itemsThisPageEnd {
				itemsThisPageEnd = nav.TotalItem
			}
			ret = ret + "<p>" + strconv.Itoa(itemsThisPageStart) + "-" + strconv.Itoa(itemsThisPageEnd) + "/" + strconv.Itoa(nav.TotalItem) + "</p>"
		}
		return template.HTML(ret)
	},
	"Sukebei":            config.IsSukebei,
	"getDefaultLanguage": languages.GetDefaultLanguage,
	"getAvatar": func(hash string, size int) string {
		return "https://www.gravatar.com/avatar/" + hash + "?s=" + strconv.Itoa(size)
	},
	"CurrentOrAdmin":       userPermission.CurrentOrAdmin,
	"CurrentUserIdentical": userPermission.CurrentUserIdentical,
	"HasAdmin":             userPermission.HasAdmin,
	"NeedsCaptcha":         userPermission.NeedsCaptcha,
	"GetRole":              userPermission.GetRole,
	"IsFollower":           userPermission.IsFollower,
	"NoEncode": func(str string) template.HTML {
		return template.HTML(str)
	},
	"calcWidthSeed": func(seed uint32, leech uint32) float64 {
		return float64(float64(seed)/(float64(seed)+float64(leech))) * 100
	},
	"calcWidthLeech": func(seed uint32, leech uint32) float64 {
		return float64(float64(leech)/(float64(seed)+float64(leech))) * 100
	},
	"formatDateRFC": func(t time.Time) string {
		// because time.* isn't available in templates...
		return t.Format(time.RFC3339)
	},
	"GetCategories": func(keepParent bool) map[string]string {
		return categories.GetCategoriesSelect(keepParent)
	},
	"CategoryName": func(category string, sub_category string) string {
		s := category + "_" + sub_category

		if category, ok := categories.GetCategories()[s]; ok {
			return category
		}
		return ""
	},
	"fileSize": func(filesize int64, T languages.TemplateTfunc) template.HTML {
		if filesize == 0 {
			return T("unknown")
		}
		return template.HTML(util.FormatFilesize(filesize))
	},
	"makeCaptchaData": func(captchaID string, T languages.TemplateTfunc) captchaData {
		return captchaData{captchaID, T}
	},
	"DefaultUserSettings": func(s string) bool {
		return config.DefaultUserSettings[s]
	},
	"makeTreeViewData": func(f *filelist.FileListFolder, nestLevel int, T languages.TemplateTfunc, identifierChain string) interface{} {
		return struct{
			Folder *filelist.FileListFolder
			NestLevel int
			T languages.TemplateTfunc
			IdentifierChain string
		}{ f, nestLevel, T, identifierChain }
	},
}
