package themes

import (
	"net/http"
	"errors"
	"github.com/NyaaPantsu/nyaa/model"
)

type UserRetriever interface {
	RetrieveCurrentUser(r *http.Request) (model.User, error)
}

var userRetriever UserRetriever = nil

func GetThemeFromRequest(r *http.Request) string {
	user, _ := getCurrentUser(r)
	if user.ID > 0 {
		return user.Theme
	}
	cookie, err := r.Cookie("theme")
	if err == nil {
		return cookie.Value
	}
	return ""
}

// Pasted from translation
func getCurrentUser(r *http.Request) (model.User, error) {
	if userRetriever == nil {
		return model.User{}, errors.New("failed to get current user: no user retriever set")
	}

	return userRetriever.RetrieveCurrentUser(r)
}
