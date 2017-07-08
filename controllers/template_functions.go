package controllers

import (
	"html/template"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/CloudyKit/jet"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/filelist"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/torrentLanguages"
)

type captchaData struct {
	CaptchaID string
	T         publicSettings.TemplateTfunc
}

// FuncMap : Functions accessible in templates by {{ $.Function }}
func templateFunctions(vars jet.VarMap) jet.VarMap {
	vars.Set("inc", func(i int) int {
		return i + 1
	})
	vars.Set("min", math.Min)
	vars.Set("genRoute", func(name string, params ...string) string {
		return "error"
	})
	vars.Set("getRawQuery", func(currentUrl *url.URL) string {
		return currentUrl.RawQuery
	})
	vars.Set("genViewTorrentRoute", func(torrent_id uint) string {
		// Helper for when you have an uint while genRoute("view_torrent", ...) takes a string
		// FIXME better solution?
		s := strconv.FormatUint(uint64(torrent_id), 10)
		url := "/view/" + s
		return url
	})
	vars.Set("genSearchWithOrdering", func(currentUrl *url.URL, sortBy string) string {
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

		u, _ := url.Parse("/search")
		u.RawQuery = values.Encode()

		return u.String()
	})
	vars.Set("genSortArrows", func(currentUrl *url.URL, sortBy string) template.HTML {
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
	})
	vars.Set("genNav", func(nav navigation, currentUrl *url.URL, pagesSelectable int) template.HTML {
		var ret = ""
		if nav.TotalItem > 0 {
			maxPages := math.Ceil(float64(nav.TotalItem) / float64(nav.MaxItemPerPage))

			if nav.CurrentPage-1 > 0 {
				url := "/" + nav.Route + "/1"
				ret = ret + "<a id=\"page-prev\" href=\"" + url + "?" + currentUrl.RawQuery + "\" aria-label=\"Previous\"><li><span aria-hidden=\"true\">&laquo;</span></li></a>"
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
				url := "/" + nav.Route + "/" + pageNum
				ret = ret + "<a aria-label=\"Page " + strconv.Itoa(i) + "\" href=\"" + url + "?" + currentUrl.RawQuery + "\">" + "<li"
				if i == nav.CurrentPage {
					ret = ret + " class=\"active\""
				}
				ret = ret + ">" + strconv.Itoa(i) + "</li></a>"
			}
			if nav.CurrentPage < int(maxPages) {
				url := "/" + nav.Route + "/" + strconv.Itoa(nav.CurrentPage+1)
				ret = ret + "<a id=\"page-next\" href=\"" + url + "?" + currentUrl.RawQuery + "\" aria-label=\"Next\"><li><span aria-hidden=\"true\">&raquo;</span></li></a>"
			}
			itemsThisPageStart := nav.MaxItemPerPage*(nav.CurrentPage-1) + 1
			itemsThisPageEnd := nav.MaxItemPerPage * nav.CurrentPage
			if nav.TotalItem < itemsThisPageEnd {
				itemsThisPageEnd = nav.TotalItem
			}
			ret = ret + "<p>" + strconv.Itoa(itemsThisPageStart) + "-" + strconv.Itoa(itemsThisPageEnd) + "/" + strconv.Itoa(nav.TotalItem) + "</p>"
		}
		return template.HTML(ret)
	})
	vars.Set("Sukebei", config.IsSukebei)
	vars.Set("getDefaultLanguage", publicSettings.GetDefaultLanguage)
	vars.Set("getAvatar", func(hash string, size int) string {
		return "https://www.gravatar.com/avatar/" + hash + "?s=" + strconv.Itoa(size)
	})

	vars.Set("DisplayTorrent", func(t models.Torrent, u *models.User) bool {
		return (!t.Hidden && t.Status != 0) || u.CurrentOrAdmin(t.UploaderID)
	})
	vars.Set("NoEncode", func(str string) template.HTML {
		return template.HTML(str)
	})
	vars.Set("calcWidthSeed", func(seed uint32, leech uint32) float64 {
		return float64(float64(seed)/(float64(seed)+float64(leech))) * 100
	})
	vars.Set("calcWidthLeech", func(seed uint32, leech uint32) float64 {
		return float64(float64(leech)/(float64(seed)+float64(leech))) * 100
	})
	vars.Set("formatDateRFC", func(t time.Time) string {
		// because time.* isn't available in templates...
		return t.Format(time.RFC3339)
	})
	vars.Set("GetHostname", format.GetHostname)
	vars.Set("GetCategories", func(keepParent bool, keepChild bool) categories.Categories {
		return categories.GetSelect(keepParent, keepChild)
	})
	vars.Set("GetCategory", func(category string, keepParent bool) (categoryRet categories.Categories) {
		cats := categories.GetSelect(true, true)
		found := false
		categoryRet = make(categories.Categories, len(cats))
		for _, v := range cats {
			if v.ID == category+"_" {
				found = true
				if keepParent {
					categoryRet = append(categoryRet, v)
				}
			} else if len(v.ID) <= 2 && len(categoryRet) > 0 {
				break
			} else if found {
				categoryRet = append(categoryRet, v)
			}
		}
		return
	})
	vars.Set("CategoryName", func(category string, sub_category string) string {
		s := category + "_" + sub_category

		if category, ok := categories.GetByID(s); ok {
			return category.Name
		}
		return ""
	})
	vars.Set("GetTorrentLanguages", torrentLanguages.GetTorrentLanguages)
	vars.Set("LanguageName", func(code string, T publicSettings.TemplateTfunc) template.HTML {
		if code == "other" || code == "multiple" {
			return T("language_" + code + "_name")
		}

		if !torrentLanguages.LanguageExists(code) {
			return T("unknown")
		}

		return T("language_" + code + "_name")
	})
	vars.Set("FlagCode", func(languageCode string) string {
		if languageCode == "other" || languageCode == "multiple" {
			return languageCode
		}

		return torrentLanguages.FlagFromLanguage(languageCode)
	})
	vars.Set("fileSize", func(filesize int64, T publicSettings.TemplateTfunc) template.HTML {
		if filesize == 0 {
			return T("unknown")
		}
		return template.HTML(format.FileSize(filesize))
	})
	vars.Set("makeCaptchaData", func(captchaID string, T publicSettings.TemplateTfunc) captchaData {
		return captchaData{captchaID, T}
	})
	vars.Set("DefaultUserSettings", func(s string) bool {
		return config.Get().Users.DefaultUserSettings[s]
	})
	vars.Set("makeTreeViewData", func(f *filelist.FileListFolder, nestLevel int, T publicSettings.TemplateTfunc, identifierChain string) interface{} {
		return struct {
			Folder          *filelist.FileListFolder
			NestLevel       int
			T               publicSettings.TemplateTfunc
			IdentifierChain string
		}{f, nestLevel, T, identifierChain}
	})
	vars.Set("lastID", func(currentUrl *url.URL, torrents []models.TorrentJSON) int {
		values := currentUrl.Query()

		order := false
		sort := "2"

		if _, ok := values["order"]; ok {
			order, _ = strconv.ParseBool(values["order"][0])
		}
		if _, ok := values["sort"]; ok {
			sort = values["sort"][0]
		}
		lastID := 0
		if sort == "2" || sort == "" {
			if order {
				lastID = int(torrents[len(torrents)-1].ID)
			} else if len(torrents) > 0 {
				lastID = int(torrents[0].ID)
			}
		}
		return lastID
	})
	vars.Set("getReportDescription", func(d string, T publicSettings.TemplateTfunc) string {
		if d == "illegal" {
			return "Illegal content"
		} else if d == "spam" {
			return "Spam / Garbage"
		} else if d == "wrongcat" {
			return "Wrong category"
		} else if d == "dup" {
			return "Duplicate / Deprecated"
		}
		return string(T(d))
	})
	vars.Set("genUploaderLink", func(uploaderID uint, uploaderName template.HTML, torrentHidden bool) template.HTML {
		uploaderID, username := torrents.HideUser(uploaderID, string(uploaderName), torrentHidden)
		if uploaderID == 0 {
			return template.HTML(username)
		}
		url := "/user/" + strconv.Itoa(int(uploaderID)) + "/" + username

		return template.HTML("<a href=\"" + url + "\">" + username + "</a>")
	})
	vars.Set("genActivityContent", func(a models.Activity, T publicSettings.TemplateTfunc) template.HTML {
		return a.ToLocale(T)
	})
	vars.Set("contains", func(arr []string, comp string) bool {
		for _, str := range arr {
			if str == comp {
				return true
			}
		}
		return false
	})
	return vars
}
