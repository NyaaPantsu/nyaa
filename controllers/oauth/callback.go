package oauth

import (
	"fmt"
	"net/url"
	"strings"

	msg "github.com/NyaaPantsu/nyaa/utils/messages"

	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
	"github.com/parnurzeal/gorequest"
	gooauth "golang.org/x/oauth2"
)

func CallbackHandler(conf gooauth.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := fosite.NewContext()
		messages := msg.GetMessages(c)
		if c.Request.URL.Query().Get("error") != "" {
			messages.AddError("errors", c.Request.URL.Query().Get("error"))
			messages.AddError("errors", c.Request.URL.Query().Get("error_description"))
			templateVariables := templates.Commonvariables(c)
			templates.Render(c, "site/api/errors.jet.html", templateVariables)
			return
		}

		if c.Request.URL.Query().Get("revoke") != "" {
			revokeURL := strings.Replace(conf.Endpoint.TokenURL, "token", "revoke", 1)
			resp, body, errs := gorequest.New().Post(revokeURL).SetBasicAuth(conf.ClientID, conf.ClientSecret).SendString(url.Values{
				"token_type_hint": {"refresh_token"},
				"token":           {c.Request.URL.Query().Get("revoke")},
			}.Encode()).End()
			if len(errs) > 0 {
				messages.AddError("errors", fmt.Sprintf(`Could not revoke token %s`, errs))
			}
			templateVariables := templates.Commonvariables(c)
			templateVariables.Set("ResponseCode", resp.StatusCode)
			if body != "" {
				templateVariables.Set("Response", fmt.Sprintf(`<p>Got a response from the revoke endpoint:<br><code>%s</code></p>`, body))
			}
			if !messages.HasErrors() {
				templateVariables.Set("Revoke", true)
			}
			templates.Render(c, "site/api/revoke.jet.html", templateVariables)
			return
		}

		if c.Request.URL.Query().Get("refresh") != "" {
			_, body, errs := gorequest.New().Post(conf.Endpoint.TokenURL).SetBasicAuth(conf.ClientID, conf.ClientSecret).SendString(url.Values{
				"grant_type":    {"refresh_token"},
				"refresh_token": {c.Request.URL.Query().Get("refresh")},
				"scope":         {"user"},
			}.Encode()).End()
			if len(errs) > 0 {
				messages.AddError("errors", fmt.Sprintf(`Could not refresh token %s`, errs))
			}
			templateVariables := templates.Commonvariables(c)
			templateVariables.Set("Response", fmt.Sprintf(`<p>Got a response from the revoke endpoint:<br><code>%s</code></p>`, body))
			if !messages.HasErrors() {
				templateVariables.Set("Response", true)
			}
			templates.Render(c, "site/api/refresh.jet.html", templateVariables)
			return
		}

		if c.Request.URL.Query().Get("code") == "" {
			messages.AddError("errors", fmt.Sprintln(`Could not find the authorize code.`))
			templateVariables := templates.Commonvariables(c)
			templates.Render(c, "site/api/errors.jet.html", templateVariables)
			return
		}

		token, err := conf.Exchange(ctx, c.Request.URL.Query().Get("code"))
		if err != nil {
			messages.AddError("errors", fmt.Sprintf(`<p>I tried to exchange the authorize code for an access token but it did not work but got error: %s</p>`, err.Error()))
			templateVariables := templates.Commonvariables(c)
			templates.Render(c, "site/api/errors.jet.html", templateVariables)
			return
		}

		messages.AddInfo("infos", fmt.Sprintf("%s", token.raw))
		templateVariables := templates.Commonvariables(c)
		templateVariables.Set("Code", c.Request.URL.Query().Get("code"))
		templateVariables.Set("AccessToken", token.AccessToken)
		templateVariables.Set("RefreshToken", token.RefreshToken)

		if !messages.HasErrors() {
			templateVariables.Set("Callback", true)
		}
		templates.Render(c, "site/api/callback.jet.html", templateVariables)
	}
}
