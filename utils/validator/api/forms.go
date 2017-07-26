package apiValidator

type CreateForm struct {
	ID                string
	Name              string   `validate:"required,min=3" form:"name"`
	RedirectURI       []string `validate:"required,dive,uri" form:"redirect_uri"`
	GrantTypes        []string `validate:"-" form:"grant_types"`
	ResponseTypes     []string `validate:"-" form:"response_types"`
	Scope             string   `validate:"-" form:"scope"`
	Owner             string   `validate:"required,min=3" form:"owner"`
	PolicyURI         string   `validate:"omitempty,uri" form:"policy_uri"`
	TermsOfServiceURI string   `validate:"omitempty,uri" form:"tos_uri"`
	ClientURI         string   `validate:"omitempty,uri" form:"client_uri"`
	LogoURI           string   `validate:"omitempty,uri" form:"logo_uri"`
	Contacts          []string `validate:"required,dive,min=3,email" form:"contacts"`
	Secret            string   `validate:"omitempty,min=8" form:"secret"`
}
