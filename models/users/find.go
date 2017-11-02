package users

import (
	"errors"
	"net/http"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/log"
)

// FindOrCreateUser creates a user.
func FindOrCreate(username string) (models.User, int, error) {
	var user models.User
	if models.ORM.Where("username=?", username).First(&user).RecordNotFound() {
		var newUser models.User
		newUser.Username = username
		log.Debugf("user %+v\n", newUser)
		if models.ORM.Create(&newUser).Error != nil {
			return newUser, http.StatusBadRequest, errors.New("user not created")
		}
		log.Debugf("retrieved User %v\n", newUser)
		return newUser, http.StatusOK, nil
	}
	return user, http.StatusBadRequest, nil
}

// RetrieveUsers retrieves users.
func GetAll() ([]*models.User, int, error) {
	var users []*models.User
	err := models.ORM.Model(&models.User{}).Find(&users).Error
	if err != nil {
		return users, http.StatusInternalServerError, err
	}
	return users, 0, nil
}

// GetByEmail retrieves a user by an email
func FindByEmail(email string) (*models.User, string, int, error) {
	var user models.User
	if models.ORM.Unscoped().Where("email = ?", email).First(&user).RecordNotFound() {
		return &user, email, http.StatusNotFound, errors.New("user_not_found")
	}
	return &user, email, http.StatusOK, nil
}

// RetrieveUserByAPIToken retrieves a user by an API token
func FindByAPIToken(apiToken string) (*models.User, string, int, error) {
	var user models.User
	if models.ORM.Unscoped().Where("api_token = ?", apiToken).First(&user).RecordNotFound() {
		return &user, apiToken, http.StatusNotFound, errors.New("user_not_found")
	}
	return &user, apiToken, http.StatusOK, nil
}

// RetrieveUserByAPITokenAndName retrieves a user by an API token and his username
func FindByAPITokenAndName(apiToken string, username string) (*models.User, string, string, int, error) {
	var user models.User
	if models.ORM.Unscoped().Where("api_token = ? AND username = ?", apiToken, username).First(&user).RecordNotFound() {
		return &user, apiToken, username, http.StatusNotFound, errors.New("user_not_found")
	}
	return &user, apiToken, username, http.StatusOK, nil
}

// RetrieveUsersByEmail retrieves users by an email
func FindUsersByEmail(email string) []*models.User {
	var users []*models.User
	models.ORM.Where("email = ?", email).Find(&users)
	return users
}

// FindByUsername retrieves a user by username.
func FindByUsername(username string) (*models.User, string, int, error) {
	var user models.User
	if models.ORM.Where("username = ?", username).First(&user).RecordNotFound() {
		return &user, username, http.StatusNotFound, errors.New("user_not_found")
	}
	return &user, username, http.StatusOK, nil
}

// RetrieveOldUploadsByUsername retrieves olduploads by username
func FindOldUploadsByUsername(username string) ([]uint, error) {
	var ret []uint
	var tmp []*models.UserUploadsOld
	err := models.ORM.Where("username = ?", username).Find(&tmp).Error
	if err != nil {
		return ret, err
	}
	for _, tmp2 := range tmp {
		ret = append(ret, tmp2.TorrentID)
	}
	return ret, nil
}

// FindByID retrieves a user by ID.
func FindByID(id uint) (*models.User, int, error) {
	var user = &models.User{}
	if models.ORM.Preload("Notifications").Last(user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("user_not_found")
	}
	var liked, likings []models.User
	models.ORM.Joins("JOIN user_follows on user_follows.user_id=?", user.ID).Where("users.user_id = user_follows.following").Group("users.user_id").Find(&likings)
	models.ORM.Joins("JOIN user_follows on user_follows.following=?", user.ID).Where("users.user_id = user_follows.user_id").Group("users.user_id").Find(&liked)
	user.Followers = likings
	user.Likings = liked
	return user, http.StatusOK, nil
}

// FindRawByID retrieves a user by ID without anything.
func FindRawByID(id uint) (*models.User, int, error) {
	var user = &models.User{}
	if models.ORM.Last(user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("user_not_found")
	}
	return user, http.StatusOK, nil
}

func SessionByID(id uint) (*models.User, int, error) {
	var user = &models.User{}
	if models.ORM.Preload("Notifications").Where("user_id = ?", id).First(user).RecordNotFound() { // We only load unread notifications
		return user, http.StatusBadRequest, errors.New("user_not_found")
	}
	return user, http.StatusOK, nil
}

// FindForAdmin retrieves a user for an administrator, preloads torrents.
func FindForAdmin(id uint) (*models.User, int, error) {
	var user = &models.User{}
	if models.ORM.Preload("Notifications").Preload("Torrents").Last(user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("user_not_found")
	}
	var liked, likings []models.User
	models.ORM.Joins("JOIN user_follows on user_follows.user_id=?", user.ID).Where("users.user_id = user_follows.following").Group("users.user_id").Find(&likings)
	models.ORM.Joins("JOIN user_follows on user_follows.following=?", user.ID).Where("users.user_id = user_follows.user_id").Group("users.user_id").Find(&liked)
	user.Followers = likings
	user.Likings = liked
	return user, http.StatusOK, nil
}

// FindUsersForAdmin retrieves users for an administrator, preloads torrents.
func FindUsersForAdmin(limit int, offset int) ([]models.User, int) {
	var users []models.User
	var nbUsers int
	models.ORM.Model(&users).Count(&nbUsers)
	models.ORM.Preload("Torrents").Limit(limit).Offset(offset).Order("user_id DESC").Find(&users)
	return users, nbUsers
}
