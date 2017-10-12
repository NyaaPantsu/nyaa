package moderatorController

import (
	"net/http"
	"strings"

	"html"
	"strconv"

	"fmt"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/oauth_client"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/oauth2/manager"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/api"
	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
)

func formClientController(c *gin.Context) {
	client := &models.OauthClient{}
	messages := msg.GetMessages(c)

	id := c.Query("id")
	if id == "" && len(messages.GetInfos("ID_TORRENT")) > 0 {
		id = messages.GetInfos("ID_TORRENT")[0]
	}
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
		RedirectURI:       strings.Split(client.RedirectURIs, "|"),
		GrantTypes:        strings.Split(client.GrantTypes, "|"),
		ResponseTypes:     strings.Split(client.ResponseTypes, "|"),
		Scope:             client.Scope,
		Owner:             client.Owner,
		PolicyURI:         client.PolicyURI,
		TermsOfServiceURI: client.TermsOfServiceURI,
		ClientURI:         client.ClientURI,
		LogoURI:           client.LogoURI,
		Contacts:          strings.Split(client.Contacts, "|"),
	}
	c.Bind(form)
	if form.ID == "" && form.Secret == "" {
		token, err := format.GenerateRandomString(32)
		if err == nil {
			form.Secret = token
		}
	}
	templates.Form(c, "admin/oauth_client_form.jet.html", form)
}

func formPostClientController(c *gin.Context) {
	messages := msg.GetMessages(c)
	sqlManager := &manager.SQLManager{&fosite.BCrypt{WorkFactor: 12}}
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
			client := manager.ToClient(form.Bind(client))
			err = sqlManager.CreateClient(client) // Making the create query through the oauth manager
			if err != nil {
				// Error, we add it as a message
				messages.AddErrorT("errors", "create_client_failed")
			} else {
				// Success, we redirect to the edit form
				messages.AddInfoT("infos", "create_client_success")
				messages.AddInfo("ID_TORRENT", client.GetID())
			}
		}
	}
	// If we are still here, we show the form
	formClientController(c)
}

// clientsListPanel : Controller for listing oauth clients, can accept pages
func clientsListPanel(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	var err error
	owner := c.Query("q")
	messages := msg.GetMessages(c)
	deleted := c.Request.URL.Query()["deleted"]

	if deleted != nil {
		messages.AddInfoTf("infos", "oauth_client_deleted")
	}
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	var conditions string
	var values []interface{}
	if owner != "" {
		conditions = "owner = ?"
		values = append(values, owner)
	}

	clients, nbClients, err := oauth_client.FindAll(offset, (pagenum-1)*offset, conditions, values)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	nav := templates.Navigation{nbClients, offset, pagenum, "mod/oauth_client/p"}
	templates.ModelList(c, "admin/clientlist.jet.html", clients, nav, templates.NewSearchForm(c))
}

// clientsDeleteModPanel : Controller for deleting a comment
func clientsDeleteModPanel(c *gin.Context) {
	id := c.PostForm("id")
	sqlManager := manager.SQLManager{&fosite.BCrypt{WorkFactor: 12}}
	client, err := oauth_client.FindByID(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	err = sqlManager.DeleteClient(id)
	if err == nil {
		activities.Log(&models.User{}, fmt.Sprintf("oauth_client_%s", client.ID), "delete", "oauth_client_deleted_by", client.ID, client.Owner, router.GetUser(c).Username)
	}

	c.Redirect(http.StatusSeeOther, "/mod/oauth_client?deleted")
}
