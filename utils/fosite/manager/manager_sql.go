package manager

import (
	"context"
	"strings"

	"net/http"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/oauth_client"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/fosite/client"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = &RichError{
		Status: http.StatusNotFound,
		error:  errors.New("Not found"),
	}
)

type SQLManager struct {
	Hasher fosite.Hasher
}

type RichError struct {
	Status int
	error
}

func sqlDataFromClient(d *client.Client) *models.OauthClient {
	return &models.OauthClient{
		ID:                d.ID,
		Name:              d.Name,
		Secret:            d.Secret,
		RedirectURIs:      strings.Join(d.RedirectURIs, "|"),
		GrantTypes:        strings.Join(d.GrantTypes, "|"),
		ResponseTypes:     strings.Join(d.ResponseTypes, "|"),
		Scope:             d.Scope,
		Owner:             d.Owner,
		PolicyURI:         d.PolicyURI,
		TermsOfServiceURI: d.TermsOfServiceURI,
		ClientURI:         d.ClientURI,
		LogoURI:           d.LogoURI,
		Contacts:          strings.Join(d.Contacts, "|"),
		Public:            d.Public,
	}
}

func (m *SQLManager) GetConcreteClient(id string) (*client.Client, error) {
	d, err := oauth_client.FindByID(id)
	if err != nil {
		return nil, err
	}
	return ToClient(d), nil
}

func (m *SQLManager) GetClient(_ context.Context, id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

func (m *SQLManager) UpdateClient(c *client.Client) error {
	o, err := m.GetClient(context.Background(), c.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	if c.Secret == "" {
		c.Secret = string(o.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash([]byte(c.Secret))
		if err != nil {
			return errors.WithStack(err)
		}
		c.Secret = string(h)
	}

	s := sqlDataFromClient(c)

	_, err = s.Update()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *SQLManager) Authenticate(id string, secret []byte) (*client.Client, error) {
	c, err := m.GetConcreteClient(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}

func (m *SQLManager) CreateClient(c *client.Client) error {
	if c.ID == "" {
		c.ID = uuid.New()
	}

	h, err := m.Hasher.Hash([]byte(c.Secret))
	if err != nil {
		return errors.WithStack(err)
	}
	c.Secret = string(h)

	data := sqlDataFromClient(c)
	err = models.ORM.Create(data).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *SQLManager) DeleteClient(id string) error {
	return oauth_client.Delete(id)
}

func (m *SQLManager) GetClients() (map[string]client.Client, error) {
	var d = []models.OauthClient{}
	clients := make(map[string]client.Client)

	err := models.ORM.Find(d).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, k := range d {
		clients[k.ID] = *ToClient(&k)
	}
	return clients, nil
}
func ToClient(d *models.OauthClient) *client.Client {
	return &client.Client{
		ID:                d.ID,
		Name:              d.Name,
		Secret:            d.Secret,
		RedirectURIs:      format.SplitNonEmpty(d.RedirectURIs, "|"),
		GrantTypes:        format.SplitNonEmpty(d.GrantTypes, "|"),
		ResponseTypes:     format.SplitNonEmpty(d.ResponseTypes, "|"),
		Scope:             d.Scope,
		Owner:             d.Owner,
		PolicyURI:         d.PolicyURI,
		TermsOfServiceURI: d.TermsOfServiceURI,
		ClientURI:         d.ClientURI,
		LogoURI:           d.LogoURI,
		Contacts:          format.SplitNonEmpty(d.Contacts, "|"),
		Public:            d.Public,
	}
}
