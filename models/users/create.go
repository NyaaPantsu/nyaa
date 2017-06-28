package users

// CreateUserFromForm creates a user from a registration form.
func CreateUserFromRequest(registrationForm *formStruct.RegistrationForm) (*models.User, error) {
	var user &models.User{}
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

	if ORM.Create(&user).Error != nil {
		return user, errors.New("user not created")
	}

	return user, nil
}

// CreateUser creates a user.
func CreateUser(c *gin.Context) int {
	var user models.User
	var registrationForm formStruct.RegistrationForm
	var status int
	var err error
	messages := msg.GetMessages(c)
	c.Bind(&registrationForm)
	usernameCandidate := SuggestUsername(registrationForm.Username)
	if usernameCandidate != registrationForm.Username {
		messages.AddErrorTf("username", "username_taken", usernameCandidate)
		return http.StatusInternalServerError
	}
	if registrationForm.Email != "" && CheckEmail(registrationForm.Email) {
		messages.AddErrorT("email", "email_in_db")
		return http.StatusInternalServerError
	}
	password, err := bcrypt.GenerateFromPassword([]byte(registrationForm.Password), 10)
	if err != nil {
		messages.ImportFromError("errors", err)
		return http.StatusInternalServerError
	}
	registrationForm.Password = string(password)
	user, err = CreateUserFromRequest(&registrationForm)
	if err != nil {
		messages.ImportFromError("errors", err)
		return http.StatusInternalServerError
	}
	if registrationForm.Email != "" {
		SendVerificationToUser(user, registrationForm.Email)
	}
	status, err = cookies.SetLoginCookies(c, user)
	if err != nil {
		messages.ImportFromError("errors", err)
	}
	return status
}