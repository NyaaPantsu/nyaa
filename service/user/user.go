package userService

import (
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	formStruct "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/util/crypto"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/ewhal/nyaa/util/timeHelper"
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

// CreateUserFromForm creates a user from a registration form.
func CreateUserFromForm(registrationForm formStruct.RegistrationForm) (model.User, error) {
	var user model.User
	log.Debugf("registrationForm %+v\n", registrationForm)
	modelHelper.AssignValue(&user, &registrationForm)
	user.Md5 = crypto.GenerateMD5Hash(user.Email)  // Gravatar
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
		return &model.PublicUser{User: &user}, isAuthor, currentUserId, http.StatusNotFound, errors.New("User is not found.")
	}
	currentUser, err := CurrentUser(r)
	if err == nil {
		currentUserId = currentUser.Id
		isAuthor = currentUser.Id == user.Id
	}

	var likings []model.User
	var likingCount int
	db.ORM.Table("users_followers").Where("users_followers.user_id=?", user.Id).Count(&likingCount)
	if err = db.ORM.Order("created_at desc").Select(config.UserPublicFields).
		Joins("JOIN users_followers on users_followers.user_id=?", user.Id).
		Where("users.id = users_followers.follower_id").
		Group("users.id").Find(&likings).Error; err != nil {
		log.Fatal(err.Error())
	}
	user.Likings = likings

	var liked []model.User
	var likedCount int
	db.ORM.Table("users_followers").Where("users_followers.follower_id=?", user.Id).Count(&likedCount)
	if err = db.ORM.Order("created_at desc").Select(config.UserPublicFields).
		Joins("JOIN users_followers on users_followers.follower_id=?", user.Id).
		Where("users.id = users_followers.user_id").
		Group("users.id").Find(&liked).Error; err != nil {
		log.Fatal(err.Error())
	}
	user.Liked = liked

	log.Debugf("user liking %v\n", user.Likings)
	log.Debugf("user liked %v\n", user.Liked)
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
	user.Md5 = crypto.GenerateMD5Hash(user.Email)
	token, err := crypto.GenerateRandomToken32()
	if err != nil {
		return http.StatusInternalServerError, errors.New("Token not generated.")
	}
	user.Token = token
	user.TokenExpiration = timeHelper.FewDaysLater(config.AuthTokenExpirationDay)
	if db.ORM.Save(user).Error != nil {
		return http.StatusInternalServerError, errors.New("User is not updated.")
	}
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

// AddRoleToUser adds a role to a user.
func AddRoleToUser(r *http.Request) (int, error) {
	var form formStruct.UserRoleForm
	var user model.User
	var role model.Role
	var roles []model.Role
	modelHelper.BindValueForm(&form, r)

	if db.ORM.First(&user, form.UserId).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	if db.ORM.First(&role, form.RoleId).RecordNotFound() {
		return http.StatusNotFound, errors.New("Role is not found.")
	}
	log.Debugf("user email : %s", user.Email)
	log.Debugf("Role name : %s", role.Name)
	db.ORM.Model(&user).Association("Roles").Append(role)
	db.ORM.Model(&user).Association("Roles").Find(&roles)
	if db.ORM.Save(&user).Error != nil {
		return http.StatusInternalServerError, errors.New("Role not appended to user.")
	}
	return http.StatusOK, nil
}

// RemoveRoleFromUser removes a role from a user.
func RemoveRoleFromUser(w http.ResponseWriter, r *http.Request, userId string, roleId string) (int, error) {
	var user model.User
	var role model.Role
	if db.ORM.First(&user, userId).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	if db.ORM.First(&role, roleId).RecordNotFound() {
		return http.StatusNotFound, errors.New("Role is not found.")
	}

	log.Debugf("user : %v\n", user)
	log.Debugf("role : %v\n", role)
	if db.ORM.Model(&user).Association("Roles").Delete(role).Error != nil {
		return http.StatusInternalServerError, errors.New("Role is not deleted from user.")
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
	db.ORM.Model(&user).Association("Languages").Find(&user.Languages)
	db.ORM.Model(&user).Association("Roles").Find(&user.Roles)
	return user, http.StatusOK, nil
}

// RetrieveUsersForAdmin retrieves users for an administrator.
func RetrieveUsersForAdmin() []model.User {
	var users []model.User
	var userArr []model.User
	db.ORM.Find(&users)
	for _, user := range users {
		db.ORM.Model(&user).Association("Languages").Find(&user.Languages)
		db.ORM.Model(&user).Association("Roles").Find(&user.Roles)
		userArr = append(userArr, user)
	}
	return userArr
}

// ActivateUser toggle activation of a user.
func ActivateUser(r *http.Request, id string) (model.User, int, error) {
	var user model.User
	var form formStruct.ActivateForm
	modelHelper.BindValueForm(&form, r)
	if db.ORM.First(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("User is not found.")
	}
	user.Activation = form.Activation
	if db.ORM.Save(&user).Error != nil {
		return user, http.StatusInternalServerError, errors.New("User not activated.")
	}
	return user, http.StatusOK, nil
}

// CreateUserAuthentication creates user authentication.
func CreateUserAuthentication(w http.ResponseWriter, r *http.Request) (int, error) {
	var form formStruct.LoginForm
	modelHelper.BindValueForm(&form, r)
	email := form.Email
	pass := form.Password
	status, err := SetCookieHandler(w, email, pass)
	return status, err
}
