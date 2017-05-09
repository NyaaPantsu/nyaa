package search

import (
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util/log"
)

type Status uint8

const (
	ShowAll Status = iota
	FilterRemakes
	Trusted
	APlus
)

type SortMode uint8

const (
	ID SortMode = iota
	Name
	Date
	Downloads
	Size
)

type Category struct {
	Main, Sub uint8
}

func (c Category) String() (s string) {
	if c.Main != 0 {
		s += strconv.Itoa(int(c.Main))
	}
	s += "_"
	if c.Sub != 0 {
		s += strconv.Itoa(int(c.Sub))
	}
	return
}

type SearchParam struct {
	Order    bool // True means acsending
	Status   Status
	Sort     SortMode
	Category Category
	Max      uint
	Query    string
}

func SearchByQuery(r *http.Request, pagenum int) (search SearchParam, tor []model.Torrents, count int, err error) {
	max, err := strconv.ParseUint(r.URL.Query().Get("max"), 10, 32)
	if err != nil {
		err = nil
		max = 50 // default Value maxPerPage
	} else if max > 300 {
		max = 300
	}
	search.Max = uint(max)

	search.Query = r.URL.Query().Get("q")

	switch s := r.URL.Query().Get("s"); s {
	case "1":
		search.Status = FilterRemakes
	case "2":
		search.Status = Trusted
	case "3":
		search.Status = APlus
	}

	catString := r.URL.Query().Get("c")
	if s := catString; len(s) > 1 && s != "_" {
		var tmp uint64
		tmp, err = strconv.ParseUint(string(s[0]), 10, 8)
		if err != nil {
			return
		}
		search.Category.Main = uint8(tmp)

		if len(s) == 3 {
			tmp, err = strconv.ParseUint(string(s[2]), 10, 8)
			if err != nil {
				return
			}
			search.Category.Sub = uint8(tmp)
		}
	}

	order_by := ""

	switch s := r.URL.Query().Get("sort"); s {
	case "1":
		search.Sort = Name
		order_by += "torrent_name"
	case "2":
		search.Sort = Date
		order_by += "date"
	case "3":
		search.Sort = Downloads
		order_by += "downloads"
	case "4":
		search.Sort = Size
		order_by += "filesize"
	default:
		order_by += "torrent_id"
	}

	order_by += " "

	switch s := r.URL.Query().Get("order"); s {
	case "true":
		search.Order = true
		order_by += "asc"
	default:
		order_by += "desc"
	}

	userId := r.URL.Query().Get("userId")

	parameters := torrentService.WhereParams{
		Params: make([]interface{}, 0, 64),
	}
	conditions := make([]string, 0, 64)
	if search.Category.Main != 0 {
		conditions = append(conditions, "category = ?")
		parameters.Params = append(parameters.Params, string(catString[0]))
	}
	if search.Category.Sub != 0 {
		conditions = append(conditions, "sub_category = ?")
		parameters.Params = append(parameters.Params, string(catString[2]))
	}
	if userId != "" {
		conditions = append(conditions, "uploader = ?")
		parameters.Params = append(parameters.Params, userId)
	}
	if search.Status != 0 {
		if search.Status == 3 {
			conditions = append(conditions, "status != ?")
		} else {
			conditions = append(conditions, "status = ?")
		}
		parameters.Params = append(parameters.Params, strconv.Itoa(int(search.Status)+1))
	}

	searchQuerySplit := strings.Fields(search.Query)
	for i, word := range searchQuerySplit {
		firstRune, _ := utf8.DecodeRuneInString(word)
		if len(word) == 1 && unicode.IsPunct(firstRune) {
			// some queries have a single punctuation character
			// which causes a full scan instead of using the index
			// and yields no meaningful results.
			// due to len() == 1 we're just looking at 1-byte/ascii
			// punctuation characters.
			continue
		}

		// TEMP: Workaround to at least make SQLite search testable for
		// development.
		// TODO: Actual case-insensitive search for SQLite
		var operator string
		if db.ORM.Dialect().GetName() == "sqlite3" {
			operator = "LIKE ?"
		} else {
			operator = "ILIKE ?"
		}

		// TODO: make this faster ?
		conditions = append(conditions, "torrent_name "+operator)
		parameters.Params = append(parameters.Params, "%"+searchQuerySplit[i]+"%")
	}

	parameters.Conditions = strings.Join(conditions[:], " AND ")
	log.Infof("SQL query is :: %s\n", parameters.Conditions)
	tor, count, err = torrentService.GetTorrentsOrderBy(&parameters, order_by, int(search.Max), int(search.Max)*(pagenum-1))
	return
}
