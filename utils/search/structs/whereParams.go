package structs

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
