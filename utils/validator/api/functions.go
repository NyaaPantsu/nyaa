package apiValidator

import (
	"strings"

	"github.com/NyaaPantsu/nyaa/models"
)

func (form *CreateForm) Bind(d *models.OauthClient) *models.OauthClient {
	d.Name = form.Name
	d.RedirectURIs = strings.Join(form.RedirectURI, "$!$")
	d.GrantTypes = strings.Join(form.GrantTypes, "$!$")
	d.ResponseTypes = strings.Join(form.ResponseTypes, "$!$")
	d.Scope = form.Scope
	d.Owner = form.Owner
	d.PolicyURI = form.PolicyURI
	d.TermsOfServiceURI = form.TermsOfServiceURI
	d.ClientURI = form.ClientURI
	d.LogoURI = form.LogoURI
	d.Contacts = strings.Join(form.Contacts, "$!$")
	d.Secret = form.Secret
	return d
}
