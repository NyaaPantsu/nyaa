package userService

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	formStruct "github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/crypto"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/modelHelper"
	"golang.org/x/crypto/bcrypt"
)

// NewCurrentUserRetriever create CurrentUserRetriever Struct for languages
func NewCurrentUserRetriever() *CurrentUserRetriever {
	return &CurrentUserRetriever{}
}

// CurrentUserRetriever struct for languages
type CurrentUserRetriever struct{}

// RetrieveCurrentUser retrieve current user for languages
func (*CurrentUserRetriever) RetrieveCurrentUser(r *http.Request) (model.User, error) {
	user, _, err := RetrieveCurrentUser(r)
	return user, err
}

// SuggestUsername suggest user's name if user's name already occupied.
func SuggestUsername(username string) string {
	var count int
	var usernameCandidate string
	db.ORM.Model(model.User{}).Where(&model.User{Username: username}).Count(&count)
	log.Debugf("count Before : %d", count)
	if count == 0 {
		return username
	}
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
	return usernameCandidate
}

// CheckEmail : check if email is in database
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
		user.MD5 = ""
	} else {
		// Despite the email not being verified yet we calculate this for convenience reasons
		var err error
		user.MD5, err = crypto.GenerateMD5Hash(user.Email)
		if err != nil {
			return user, err
		}
	}
	user.Email = "" // unset email because it will be verified later
	user.CreatedAt = time.Now()
	// User settings to default
	user.Settings.ToDefault()
	user.SaveSettings()
	// currently unused but needs to be set:
	user.APIToken, _ = crypto.GenerateRandomToken32()
	user.APITokenExpiry = time.Unix(0, 0)

	if db.ORM.Create(&user).Error != nil {
		return user, errors.New("user not created")
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
	if usernameCandidate != registrationForm.Username {
		return http.StatusInternalServerError, fmt.Errorf("Username already taken, you can choose: %s", usernameCandidate)
	}
	if registrationForm.Email != "" && CheckEmail(registrationForm.Email) {
		return http.StatusInternalServerError, errors.New("email address already in database")
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
	if registrationForm.Email != "" {
		SendVerificationToUser(user, registrationForm.Email)
	}
	status, err = RegisterHandler(w, r)
	return status, err
}

// RetrieveUser retrieves a user.
func RetrieveUser(r *http.Request, id string) (*model.PublicUser, bool, uint, int, error) {
	var user model.User
	var currentUserID uint
	var isAuthor bool

	if db.ORM.First(&user, id).RecordNotFound() {
		return nil, isAuthor, currentUserID, http.StatusNotFound, errors.New("user not found")
	}
	currentUser, err := CurrentUser(r)
	if err == nil {
		currentUserID = currentUser.ID
		isAuthor = currentUser.ID == user.ID
	}

	return &model.PublicUser{User: &user}, isAuthor, currentUserID, http.StatusOK, nil
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
		user.MD5 = ""
	} else {
		var err error
		user.MD5, err = crypto.GenerateMD5Hash(user.Email)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	user.UpdatedAt = time.Now()
	err := db.ORM.Save(user).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// UpdateUser updates a user.
func UpdateUser(w http.ResponseWriter, form *formStruct.UserForm, formSet *formStruct.UserSettingsForm, currentUser *model.User, id string) (model.User, int, error) {
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("user not found")
	}
	log.Infof("updateUser")
	if form.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.CurrentPassword))
		if err != nil && !userPermission.HasAdmin(currentUser) {
			log.Error("Password Incorrect.")
			return user, http.StatusInternalServerError, errors.New("password incorrect")
		}
		newPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
		if err != nil {
			return user, http.StatusInternalServerError, errors.New("password not generated")
		}
		form.Password = string(newPassword)
	} else { // Then no change of password
		form.Password = user.Password
	}
	if !userPermission.HasAdmin(currentUser) { // We don't want users to be able to modify some fields
		form.Status = user.Status
		form.Username = user.Username
	}
	if form.Email != user.Email {
		// send verification to new email and keep old
		SendVerificationToUser(user, form.Email)
		form.Email = user.Email
	}
	log.Debugf("form %+v\n", form)
	modelHelper.AssignValue(&user, form)

	// We set settings here
	user.ParseSettings()
	user.Settings.Set("new_torrent", formSet.NewTorrent)
	user.Settings.Set("new_torrent_email", formSet.NewTorrentEmail)
	user.Settings.Set("new_comment", formSet.NewComment)
	user.Settings.Set("new_comment_email", formSet.NewCommentEmail)
	user.Settings.Set("new_responses", formSet.NewResponses)
	user.Settings.Set("new_responses_email", formSet.NewResponsesEmail)
	user.Settings.Set("new_follower", formSet.NewFollower)
	user.Settings.Set("new_follower_email", formSet.NewFollowerEmail)
	user.Settings.Set("followed", formSet.Followed)
	user.Settings.Set("followed_email", formSet.FollowedEmail)
	user.SaveSettings()

	status, err := UpdateUserCore(&user)
	return user, status, err
}

// DeleteUser deletes a user.
func DeleteUser(w http.ResponseWriter, currentUser *model.User, id string) (int, error) {
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("user not found")
	}
	if user.ID == 0 {
		return http.StatusInternalServerError, errors.New("You can't delete that")
	}
	if db.ORM.Delete(&user).Error != nil {
		return http.StatusInternalServerError, errors.New("user not deleted")
	}
	if userPermission.CurrentUserIdentical(currentUser, user.ID) {
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
		return &model.PublicUser{User: &user}, email, http.StatusNotFound, errors.New("user not found")
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
		return &model.PublicUser{User: &user}, username, http.StatusNotFound, errors.New("user not found")
	}
	return &model.PublicUser{User: &user}, username, http.StatusOK, nil
}

// RetrieveOldUploadsByUsername retrieves olduploads by username
func RetrieveOldUploadsByUsername(username string) ([]uint, error) {
	var ret []uint
	var tmp []*model.UserUploadsOld
	err := db.ORM.Where("username = ?", username).Find(&tmp).Error
	if err != nil {
		return ret, err
	}
	for _, tmp2 := range tmp {
		ret = append(ret, tmp2.TorrentID)
	}
	return ret, nil
}

// RetrieveUserForAdmin retrieves a user for an administrator.
func RetrieveUserForAdmin(id string) (model.User, int, error) {
	var user model.User
	if db.ORM.Preload("Notifications").Preload("Torrents").Last(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("user not found")
	}
	var liked, likings []model.User
	db.ORM.Joins("JOIN user_follows on user_follows.user_id=?", user.ID).Where("users.user_id = user_follows.following").Group("users.user_id").Find(&likings)
	db.ORM.Joins("JOIN user_follows on user_follows.following=?", user.ID).Where("users.user_id = user_follows.user_id").Group("users.user_id").Find(&liked)
	user.Likings = likings
	user.Liked = liked
	return user, http.StatusOK, nil
}

// RetrieveUsersForAdmin retrieves users for an administrator.
func RetrieveUsersForAdmin(limit int, offset int) ([]model.User, int) {
	var users []model.User
	var nbUsers int
	db.ORM.Model(&users).Count(&nbUsers)
	db.ORM.Preload("Torrents").Limit(limit).Offset(offset).Find(&users)
	return users, nbUsers
}

// GetLiked : Gets who is following the user
func GetLiked(user *model.User) *model.User {
	var liked []model.User
	db.ORM.Joins("JOIN user_follows on user_follows.following=?", user.ID).Where("users.user_id = user_follows.user_id").Group("users.user_id").Find(&liked)
	user.Liked = liked
	return user
}

// GetLikings : Gets who is followed by the user
func GetLikings(user *model.User) *model.User {
	var likings []model.User
	db.ORM.Joins("JOIN user_follows on user_follows.user_id=?", user.ID).Where("users.user_id = user_follows.following").Group("users.user_id").Find(&likings)
	user.Likings = likings
	return user
}

// CreateUserAuthentication creates user authentication.
func CreateUserAuthentication(w http.ResponseWriter, r *http.Request) (int, error) {
	var form formStruct.LoginForm
	modelHelper.BindValueForm(&form, r)
	username := form.Username
	pass := form.Password
	status, err := SetCookieHandler(w, r, username, pass)
	return status, err
}

// SetFollow : Makes a user follow another
func SetFollow(user *model.User, follower *model.User) {
	if follower.ID > 0 && user.ID > 0 {
		var userFollows = model.UserFollows{UserID: user.ID, FollowerID: follower.ID}
		db.ORM.Create(&userFollows)
	}
}

// RemoveFollow : Remove a user following another
func RemoveFollow(user *model.User, follower *model.User) {
	if follower.ID > 0 && user.ID > 0 {
		var userFollows = model.UserFollows{UserID: user.ID, FollowerID: follower.ID}
		db.ORM.Delete(&userFollows)
	}
}
