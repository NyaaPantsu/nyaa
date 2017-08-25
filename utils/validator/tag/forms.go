package tagsValidator

type CreateForm struct {
	Tag  string `validate:"required" form:"tag" json:"tag"`
	Type string `validate:"required" form:"type" json:"type"`
}
