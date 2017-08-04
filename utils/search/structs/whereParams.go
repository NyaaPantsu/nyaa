package structs

import (
	"fmt"
)

// WhereParams struct for search
type WhereParams struct {
	Conditions string // Ex : name LIKE ? AND category_id LIKE ?
	Params     []interface{}
}

// CreateWhereParams : function to create WhereParams struct for search
func CreateWhereParams(conditions string, params ...interface{}) WhereParams {
	whereParams := WhereParams{
		Conditions: conditions,
		Params:     make([]interface{}, len(params)),
	}
	copy(params, whereParams.Params)

	return whereParams
}

// Identifier returns an unique identifier for whereparams
func (w *WhereParams) Identifier() string {
	params := ""
	for _, param := range w.Params {
		params += fmt.Sprintf("%v", param)
	}
	return params
}
