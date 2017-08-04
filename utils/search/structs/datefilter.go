package structs

import (
	"strconv"
	"time"
)

// DateFilter date to filter for
type DateFilter string

// ParseOld parses a date to a datefilter object and return true if it succeeded
// Deprecated: This function will be removed in 1.2, please use Parse
func (d *DateFilter) ParseOld(s string, dateType string) bool {
	if s == "" {
		*d = ""
		return false
	}
	dateInt, err := strconv.Atoi(s)
	if err != nil {
		*d = ""
		return false
	}
	switch dateType {
	case "m":
		*d = DateFilter(time.Now().AddDate(0, -dateInt, 0).Format("2006-01-02"))
	case "y":
		*d = DateFilter(time.Now().AddDate(-dateInt, 0, 0).Format("2006-01-02"))
	default:
		*d = DateFilter(time.Now().AddDate(0, 0, -dateInt).Format("2006-01-02"))
	}
	return true
}

// Parse parses a date to a datefilter object and return true if it succeeded
// This functions accept only date formatted in this way YYYY/MM/DD
func (d *DateFilter) Parse(s string) bool {
	if s == "" {
		*d = ""
		return false
	}
	date, err := time.Parse("2006/01/02", s)
	if err != nil {
		*d = ""
		return false
	}
	*d = DateFilter(date.Format("2006-01-02"))
	return true
}

func backwardCompatibility(max string, from string, to string, dtype string) (DateFilter, DateFilter) {
	fromDate, toDate := DateFilter(""), DateFilter("")
	maxage, err := strconv.Atoi(max)
	if err != nil {
		// if we can't convert it, this means maxage was not provided or wrongly, so we try to filter the query with other date filter args
		// Deprecated : We have to give backward compatibility here to dateType
		// This will be removed on 1.2, please remove from using it and go to maxage or Date argument
		if dtype != "" {
			// if to xxx is not provided, fromDate is equal to from xxx
			if to != "" {
				fromDate.ParseOld(to, dtype)
				toDate.ParseOld(from, dtype)
				return fromDate, toDate
			}

			fromDate.ParseOld(from, dtype)
			return fromDate, toDate
		}
		// This will be the future default behavior
		// We only try to parse the dates sent and if it works we assign them
		fromDate.Parse(from)
		toDate.Parse(to)
		return fromDate, toDate
	}
	// Maxage behavior where we convert the substracted date to string
	fromDate = DateFilter(time.Now().AddDate(0, 0, -maxage).Format("2006-01-02"))
	return fromDate, toDate
}
