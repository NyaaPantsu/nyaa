package router

import (
	"html/template"
	"log"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/nicksnyder/go-i18n/i18n"
)

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
		order := false  //Default is DESC
		sort := "2" //Default is Date (Actually ID, but Date is the same thing)

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
	"genNav": func(nav Navigation, currentUrl *url.URL, pagesSelectable int) template.HTML {
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
		}
		return template.HTML(ret)
	},
	"Sukebei": func() bool {
		if config.TableName == "sukebei_torrents" {
			return true
		} else {
			return false
		}
	},
	"T":                  i18n.IdentityTfunc,
	"Ts":                 i18n.IdentityTfunc,
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
    "Category_Sukebei": func(category string, sub_category string) string {
        s := category + "_" + sub_category; e := ""
        switch s {
            default:        e = ""
            case "1_":      e = "art"
            case "1_1":     e = "art_anime"
            case "1_2":     e = "art_doujinshi"
            case "1_3":     e = "art_games"
            case "1_4":     e = "art_manga"
            case "1_5":     e = "art_pictures"
            case "2_":      e = "real_life"
            case "2_1":     e = "real_life_photobooks_and_pictures"
            case "2_2":     e = "real_life_videos"
        }
        return e
    },
    "Category_Nyaa": func(category string, sub_category string) string {
        s := category + "_" + sub_category; e := ""
        switch s {
            default:        e = ""
            case "3_":      e = "anime"
            case "3_12":    e = "anime_amv"
            case "3_5":     e = "anime_english_translated"
            case "3_13":    e = "anime_non_english_translated"
            case "3_6":     e = "anime_raw"
            case "2_":      e = "audio"
            case "2_3":     e = "audio_lossless"
            case "2_4":     e = "audio_lossy"
            case "4_":      e = "literature"
            case "4_7":     e = "literature_english_translated"
            case "4_8":     e = "literature_raw"
            case "4_14":    e = "literature_non_english_translated"
            case "5_":      e = "live_action"
            case "5_9":     e = "live_action_english_translated"
            case "5_10":    e = "live_action_idol_pv"
            case "5_18":    e = "live_action_non_english_translated"
            case "5_11":    e = "live_action_raw"
            case "6_":      e = "pictures"
            case "6_15":    e = "pictures_graphics"
            case "6_16":    e = "pictures_photos"
            case "1_":      e = "software"
            case "1_1":     e = "software_applications"
            case "1_2":     e = "software_games"
        }
        return e
    },
}
