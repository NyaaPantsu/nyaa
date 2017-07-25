package moderatorController

import (
	"net/http"
	"strings"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/oauth_client"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/fosite/manager"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/api"
	"github.com/gin-gonic/gin"
)

func formClientController(c *gin.Context) {
	client := &models.OauthClient{}
	id := c.Query("id")
	if id != "" {
		var err error
		client, err = oauth_client.FindByID(id)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}

	form := &apiValidator.CreateForm{
		ID:                client.ID,
		Name:              client.Name,
		RedirectURI:       strings.Split(client.RedirectURIs, "$!$"),
		GrantTypes:        strings.Split(client.GrantTypes, "$!$"),
		ResponseTypes:     strings.Split(client.ResponseTypes, "$!$"),
		Scope:             client.Scope,
		Owner:             client.Owner,
		PolicyURI:         client.PolicyURI,
		TermsOfServiceURI: client.TermsOfServiceURI,
		ClientURI:         client.ClientURI,
		LogoURI:           client.LogoURI,
		Contacts:          strings.Split(client.Contacts, "$!$"),
	}
	c.Bind(form)

	templates.Form(c, "admin/clientapi.jet.html", form)
}

func formPostClientController(c *gin.Context) {
	messages := msg.GetMessages(c)
	sqlManager := &manager.SQLManager{}
	client := &models.OauthClient{}
	id := c.Query("id")
	if id != "" {
		var err error
		client, err = oauth_client.FindByID(id)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
	form := &apiValidator.CreateForm{}
	// We bind the request to the form
	c.Bind(form)
	// We try to validate the form
	validator.ValidateForm(form, messages)
	// If validation has failed, errors are added in messages variable
	if !messages.HasErrors() {
		// No errors, check if we update or create
		if id != "" { // Client exists we update
			err := sqlManager.UpdateClient(manager.ToClient(form.Bind(client))) // Making the update query through the oauth manager
			if err != nil {
				// Error, we add it to the messages variable
				messages.AddErrorT("errors", "update_client_failed")
			} else {
				// Success, we add a notice to the messages variable
				messages.AddInfoT("infos", "update_client_success")
			}
		} else { // Client doesn't exist, we create it
			var err error
			err = sqlManager.CreateClient(manager.ToClient(form.Bind(client))) // Making the create query through the oauth manager
			if err != nil {
				// Error, we add it as a message
				messages.AddErrorT("errors", "create_client_failed")
			} else {
				// Success, we redirect to the edit form
				messages.AddInfoT("infos", "create_client_success")
			}
		}
	}
	// If we are still here, we show the form
	formClientController(c)
}
