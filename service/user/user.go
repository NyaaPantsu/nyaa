package userService

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	formStruct "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/util/crypto"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/ewhal/nyaa/util/timeHelper"
	"fmt"
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
		return true
	}
	var count int
	db.ORM.Model(model.User{}).Where("email = ?", email).Count(&count)
	if count == 0 {
		return false // duplicate
	}
	return true
}

// CreateUserFromForm creates a user from a registration form.
func CreateUserFromForm(registrationForm formStruct.RegistrationForm) (model.User, error) {
	var user model.User
	log.Debugf("registrationForm %+v\n", registrationForm)
	modelHelper.AssignValue(&user, &registrationForm)
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
	if (usernameCandidate != registrationForm.Username) {
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
	if db.ORM.Select(config.UserPublicFields).First(&user, id).RecordNotFound() {
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
	db.ORM.Select(config.UserPublicFields).Find(&users)
	for _, user := range users {
		userArr = append(userArr, &model.PublicUser{User: user})
	}
	return userArr
}

// UpdateUserCore updates a user. (Applying the modifed data of user).
func UpdateUserCore(user *model.User) (int, error) {
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
func UpdateUser(w http.ResponseWriter, r *http.Request, id string) (*model.User, int, error) {
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return &user, http.StatusNotFound, errors.New("User is not found.")
	}
	switch r.FormValue("type") {
	case "password":
		var passwordForm formStruct.PasswordForm
		modelHelper.BindValueForm(&passwordForm, r)
		log.Debugf("form %+v\n", passwordForm)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordForm.CurrentPassword))
		if err != nil {
			log.Error("Password Incorrect.")
			return &user, http.StatusInternalServerError, errors.New("User is not updated. Password Incorrect.")
		} else {
			newPassword, err := bcrypt.GenerateFromPassword([]byte(passwordForm.Password), 10)
			if err != nil {
				return &user, http.StatusInternalServerError, errors.New("User is not updated. Password not Generated.")
			} else {
				passwordForm.Password = string(newPassword)
				modelHelper.AssignValue(&user, &passwordForm)
			}
		}
	default:
		var form formStruct.UserForm
		modelHelper.BindValueForm(&form, r)
		log.Debugf("form %+v\n", form)
		modelHelper.AssignValue(&user, &form)
	}

	status, err := UpdateUserCore(&user)
	if err != nil {
		return &user, status, err
	}
	status, err = SetCookie(w, user.Token)
	return &user, status, err
}

// DeleteUser deletes a user.
func DeleteUser(w http.ResponseWriter, id string) (int, error) {
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	if db.ORM.Delete(&user).Error != nil {
		return http.StatusInternalServerError, errors.New("User is not deleted.")
	}
	status, err := ClearCookie(w)
	return status, err
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
	if db.ORM.Unscoped().Select(config.UserPublicFields).Where("email like ?", "%"+email+"%").First(&user).RecordNotFound() {
		return &model.PublicUser{User: &user}, email, http.StatusNotFound, errors.New("User is not found.")
	}
	return &model.PublicUser{User: &user}, email, http.StatusOK, nil
}

// RetrieveUsersByEmail retrieves users by an email
func RetrieveUsersByEmail(email string) []*model.PublicUser {
	var users []*model.User
	var userArr []*model.PublicUser
	db.ORM.Select(config.UserPublicFields).Where("email like ?", "%"+email+"%").Find(&users)
	for _, user := range users {
		userArr = append(userArr, &model.PublicUser{User: user})
	}
	return userArr
}

// RetrieveUserByUsername retrieves a user by username.
func RetrieveUserByUsername(username string) (*model.PublicUser, string, int, error) {
	var user model.User
	if db.ORM.Unscoped().Select(config.UserPublicFields).Where("username like ?", "%"+username+"%").First(&user).RecordNotFound() {
		return &model.PublicUser{User: &user}, username, http.StatusNotFound, errors.New("User is not found.")
	}
	return &model.PublicUser{User: &user}, username, http.StatusOK, nil
}

// RetrieveUserForAdmin retrieves a user for an administrator.
func RetrieveUserForAdmin(id string) (model.User, int, error) {
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("User is not found.")
	}
	db.ORM.Model(&user)
	return user, http.StatusOK, nil
}

// RetrieveUsersForAdmin retrieves users for an administrator.
func RetrieveUsersForAdmin() []model.User {
	var users []model.User
	var userArr []model.User
	db.ORM.Find(&users)
	for _, user := range users {
		db.ORM.Model(&user)
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
