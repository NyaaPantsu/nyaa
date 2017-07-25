package apiValidator

type CreateForm struct {
	ID                string
	Name              string   `validate:"required,min=3"`
	RedirectURI       []string `validate:"required,dive,uri" form:"redirect_uri"`
	GrantTypes        []string `form:"grant_types"`
	ResponseTypes     []string `form:"response_types"`
	Scope             string   `form:"scope"`
	Owner             string   `validate:"required,min=3"`
	PolicyURI         string   `validate:"uri" form:"policy_uri"`
	TermsOfServiceURI string   `validate:"uri" form:"tos_uri"`
	ClientURI         string   `validate:"uri" form:"client_uri"`
	LogoURI           string   `validate:"uri" form:"logo_uri"`
	Contacts          []string `validate:"required,dive,min=3,email"`
	Secret            string   `validate:"min=8"`
}
