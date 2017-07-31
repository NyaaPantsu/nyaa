package tagsValidator

type CreateForm struct {
	Tag  string `validate:"required" form:"tag"`
	Type string `validate:"required" form:"type"`
}
