package userService

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	formStruct "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util/crypto"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/ewhal/nyaa/util/timeHelper"
	"golang.org/x/crypto/bcrypt"
)

var userFields []string = []string{"name", "email", "createdAt", "updatedAt"}

// SuggestUsername suggest user's name if user's name already occupied.
func SuggestUsername(username string) string {
	var count int
	var usernameCandidate string
	db.ORM.Model(model.User{}).Where(&model.User{Username: username}).Count(&count)
	log.Debugf("count Before : %d", count)
	if count == 0 {
		return username
	} else {
		var postfix int
		for {
			usernameCandidate = username + strconv.Itoa(postfix)
			log.Debugf("usernameCandidate: %s\n", usernameCandidate)
			db.ORM.Model(model.User{}).Where(&model.User{Username: usernameCandidate}).Count(&count)
			log.Debugf("count after : %d\n", count)
			postfix = postfix + 1
			if count == 0 {
				break
			}
		}
	}
	return usernameCandidate
}

func CheckEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	var count int
	db.ORM.Model(model.User{}).Where("email = ?", email).Count(&count)
	if count != 0 {
		return true // error: duplicate
	}
	return false
}

// CreateUserFromForm creates a user from a registration form.
func CreateUserFromForm(registrationForm formStruct.RegistrationForm) (model.User, error) {
	var user model.User
	log.Debugf("registrationForm %+v\n", registrationForm)
	modelHelper.AssignValue(&user, &registrationForm)
	if user.Email == "" {
		user.Md5 = ""
	} else {
		user.Md5 = crypto.GenerateMD5Hash(user.Email)
	}
	token, err := crypto.GenerateRandomToken32()
	if err != nil {
		return user, errors.New("Token not generated.")
	}

	user.Token = token
	user.TokenExpiration = timeHelper.FewDaysLater(config.AuthTokenExpirationDay)
	log.Debugf("user %+v\n", user)
	if db.ORM.Create(&user).Error != nil {
		return user, errors.New("User is not created.")
	}

	user.CreatedAt = time.Now()
	return user, nil
}

// CreateUser creates a user.
func CreateUser(w http.ResponseWriter, r *http.Request) (int, error) {
	var user model.User
	var registrationForm formStruct.RegistrationForm
	var status int
	var err error

	modelHelper.BindValueForm(&registrationForm, r)
	usernameCandidate := SuggestUsername(registrationForm.Username)
	if usernameCandidate != registrationForm.Username {
		return http.StatusInternalServerError, fmt.Errorf("Username already taken, you can choose: %s", usernameCandidate)
	}
	if CheckEmail(registrationForm.Email) {
		return http.StatusInternalServerError, errors.New("Email Address already in our database!")
	}
	password, err := bcrypt.GenerateFromPassword([]byte(registrationForm.Password), 10)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	registrationForm.Password = string(password)
	user, err = CreateUserFromForm(registrationForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	SendVerificationToUser(user)
	status, err = RegisterHandler(w, r)
	return status, err
}

// RetrieveUser retrieves a user.
func RetrieveUser(r *http.Request, id string) (*model.PublicUser, bool, uint, int, error) {
	var user model.User
	var currentUserId uint
	var isAuthor bool
	// var publicUser *model.PublicUser
	// publicUser.User = &user
	if db.ORM.First(&user, id).RecordNotFound() {
		return nil, isAuthor, currentUserId, http.StatusNotFound, errors.New("User is not found.")
	}
	currentUser, err := CurrentUser(r)
	if err == nil {
		currentUserId = currentUser.Id
		isAuthor = currentUser.Id == user.Id
	}

	return &model.PublicUser{User: &user}, isAuthor, currentUserId, http.StatusOK, nil
}

// RetrieveUsers retrieves users.
func RetrieveUsers() []*model.PublicUser {
	var users []*model.User
	var userArr []*model.PublicUser
	for _, user := range users {
		userArr = append(userArr, &model.PublicUser{User: user})
	}
	return userArr
}

// UpdateUserCore updates a user. (Applying the modifed data of user).
func UpdateUserCore(user *model.User) (int, error) {
	if user.Email == "" {
		user.Md5 = ""
	} else {
		user.Md5 = crypto.GenerateMD5Hash(user.Email)
	}

	token, err := crypto.GenerateRandomToken32()
	if err != nil {
		return http.StatusInternalServerError, errors.New("Token not generated.")
	}
	user.Token = token
	user.TokenExpiration = timeHelper.FewDaysLater(config.AuthTokenExpirationDay)
	if db.ORM.Save(user).Error != nil {
		return http.StatusInternalServerError, errors.New("User is not updated.")
	}

	user.UpdatedAt = time.Now()
	return http.StatusOK, nil
}

// UpdateUser updates a user.
func UpdateUser(w http.ResponseWriter, form *formStruct.UserForm, currentUser *model.User, id string) (model.User, int, error) {
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("User is not found.")
	}

	if form.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.CurrentPassword))
		if err != nil && !userPermission.HasAdmin(currentUser) {
			log.Error("Password Incorrect.")
			return user, http.StatusInternalServerError, errors.New("User is not updated. Password Incorrect.")
		} else {
			newPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
			if err != nil {
				return user, http.StatusInternalServerError, errors.New("User is not updated. Password not Generated.")
			} else {
				form.Password = string(newPassword)
			}
		}
	} else { // Then no change of password
		form.Password = user.Password
	}
	if !userPermission.HasAdmin(currentUser) { // We don't want users to be able to modify some fields
		form.Status = user.Status
		form.Username = user.Username
	}
	log.Debugf("form %+v\n", form)
	modelHelper.AssignValue(&user, form)

	status, err := UpdateUserCore(&user)
	if err != nil {
		return user, status, err
	}
	if userPermission.CurrentUserIdentical(currentUser, user.Id) {
		status, err = SetCookie(w, user.Token)
	}
	return user, status, err
}

// DeleteUser deletes a user.
func DeleteUser(w http.ResponseWriter, currentUser *model.User, id string) (int, error) {
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	if db.ORM.Delete(&user).Error != nil {
		return http.StatusInternalServerError, errors.New("User is not deleted.")
	}
	if userPermission.CurrentUserIdentical(currentUser, user.Id) {
		return ClearCookie(w)
	}
	return http.StatusOK, nil
}

// RetrieveCurrentUser retrieves a current user.
func RetrieveCurrentUser(r *http.Request) (model.User, int, error) {
	user, err := CurrentUser(r)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

// RetrieveUserByEmail retrieves a user by an email
func RetrieveUserByEmail(email string) (*model.PublicUser, string, int, error) {
	var user model.User
	if db.ORM.Unscoped().Where("email = ?", email).First(&user).RecordNotFound() {
		return &model.PublicUser{User: &user}, email, http.StatusNotFound, errors.New("User is not found.")
	}
	return &model.PublicUser{User: &user}, email, http.StatusOK, nil
}

// RetrieveUsersByEmail retrieves users by an email
func RetrieveUsersByEmail(email string) []*model.PublicUser {
	var users []*model.User
	var userArr []*model.PublicUser
	db.ORM.Where("email = ?", email).Find(&users)
	for _, user := range users {
		userArr = append(userArr, &model.PublicUser{User: user})
	}
	return userArr
}

// RetrieveUserByUsername retrieves a user by username.
func RetrieveUserByUsername(username string) (*model.PublicUser, string, int, error) {
	var user model.User
	if db.ORM.Where("username = ?", username).First(&user).RecordNotFound() {
		return &model.PublicUser{User: &user}, username, http.StatusNotFound, errors.New("User is not found.")
	}
	return &model.PublicUser{User: &user}, username, http.StatusOK, nil
}

// RetrieveUserForAdmin retrieves a user for an administrator.
func RetrieveUserForAdmin(id string) (model.User, int, error) {
	var user model.User
	if db.ORM.Preload("Torrents").First(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("User is not found.")
	}
	var liked,likings []model.User 
	db.ORM.Joins("JOIN user_follows on user_follows.user_id=?", user.Id).Where("users.user_id = user_follows.following").Group("users.user_id").Find(&likings)
	db.ORM.Joins("JOIN user_follows on user_follows.following=?", user.Id).Where("users.user_id = user_follows.user_id").Group("users.user_id").Find(&liked)
	user.Likings = likings
	user.Liked = liked
	return user, http.StatusOK, nil
}

// RetrieveUsersForAdmin retrieves users for an administrator.
func RetrieveUsersForAdmin() []model.User {
	var users []model.User
	var userArr []model.User
	db.ORM.Find(&users)
	for _, user := range users {
		db.ORM.Model(&user)
		db.ORM.Model(&user).Related("Torrents").Related("Likings").Related("Liked").Find(&model.Torrents{})
		userArr = append(userArr, user)
	}
	return userArr
}

// CreateUserAuthentication creates user authentication.
func CreateUserAuthentication(w http.ResponseWriter, r *http.Request) (int, error) {
	var form formStruct.LoginForm
	modelHelper.BindValueForm(&form, r)
	username := form.Username
	pass := form.Password
	status, err := SetCookieHandler(w, username, pass)
	return status, err
}

func SetFollow(user *model.User, follower *model.User) {
	if (follower.Id > 0 && user.Id > 0) {
		var userFollows = model.UserFollows{UserID: user.Id, FollowerID: follower.Id}
		db.ORM.Create(&userFollows)
	}
}

func RemoveFollow(user *model.User, follower *model.User) {
	if (follower.Id > 0 && user.Id > 0) {
		var userFollows = model.UserFollows{UserID: user.Id, FollowerID: follower.Id}
		db.ORM.Delete(&userFollows)
	}
}