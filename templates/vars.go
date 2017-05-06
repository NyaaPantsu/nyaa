package templates

/*
 * Variables used by the upper ones
 */
type Navigation struct {
	TotalItem      int
	MaxItemPerPage int
	CurrentPage    int
	Route          string
}

type SearchForm struct {
	Query              string
	Status             string
	Category           string
	Sort               string
	Order              string
	HideAdvancedSearch bool
}

// Some Default Values to ease things out
func NewSearchForm(params ...string) (searchForm SearchForm) {
	if len(params) > 1 {
		searchForm.Category = params[0]
	} else {
		searchForm.Category = "_"
	}
	if len(params) > 2 {
		searchForm.Sort = params[1]
	} else {
		searchForm.Sort = "torrent_id"
	}
	if len(params) > 3 {
		order := params[2]
		if order == "DESC" {
			searchForm.Order = order
		} else if order == "ASC" {
			searchForm.Order = order
		} else {
			// TODO: handle invalid value (?)
		}
	} else {
		searchForm.Order = "DESC"
	}
	return
}
