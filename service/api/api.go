package apiService

import (
	"reflect"

	"github.com/ewhal/nyaa/service/torrent"
)

type torrentsQuery struct {
	Category    int `json:"category"`
	SubCategory int `json:"sub_category"`
	Status      int `json:"status"`
	Uploader    int `json:"uploader"`
	Downloads   int `json:"downloads"`
}

type TorrentsRequest struct {
	Query      torrentsQuery `json:"search"`
	Page       int           `json:"page"`
	MaxPerPage int           `json:"limit"`
}

func (r *TorrentsRequest) ToParams() torrentService.WhereParams {
	res := torrentService.WhereParams{}
	conditions := ""
	v := reflect.ValueOf(r.Query)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			if i != 0 {
				conditions += " AND "
			}
			conditions += v.Type().Field(i).Tag.Get("json") + " = ?"
			res.Params = append(res.Params, field.Interface())
		}
	}
	res.Conditions = conditions
	return res
}
