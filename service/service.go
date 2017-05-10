package serviceBase

type WhereParams struct {
	Conditions string // Ex : name LIKE ? AND category_id LIKE ?
	Params     []interface{}
}

func CreateWhereParams(conditions string, params ...string) WhereParams {
	whereParams := WhereParams{
		Conditions: conditions,
		Params:     make([]interface{}, len(params)),
	}
	for i := range params {
		whereParams.Params[i] = params[i]
	}
	return whereParams
}
