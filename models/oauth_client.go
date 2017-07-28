package models

import (
	"errors"
	"net/http"

	"github.com/fatih/structs"
)

type OauthClient struct {
	ID                string `gorm:"column:id;primary_key;not null"`
	Name              string `gorm:"column:client_name;not null"`
	Secret            string `gorm:"column:client_secret;not null"`
	RedirectURIs      string `gorm:"column:redirect_uris;not null"`
	GrantTypes        string `gorm:"column:grant_types;not null"`
	ResponseTypes     string `gorm:"column:response_types;not null"`
	Scope             string `gorm:"column:scope;not null"`
	Owner             string `gorm:"column:owner;not null"`
	PolicyURI         string `gorm:"column:policy_uri;not null"`
	TermsOfServiceURI string `gorm:"column:tos_uri;not null"`
	ClientURI         string `gorm:"column:client_uri;not null"`
	LogoURI           string `gorm:"column:logo_uri;not null"`
	Contacts          string `gorm:"column:contacts;not null"`
	Public            bool   `gorm:"column:public;not null"`
}

func (d OauthClient) TableName() string {
	return "hydra_client"
}

func (d *OauthClient) Update() (int, error) {
	if ORM.Model(d).UpdateColumn(d.toMap()).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}
	return http.StatusOK, nil
}

// toMap : convert the model to a map of interface
func (d *OauthClient) toMap() map[string]interface{} {
	return structs.Map(d)
}
