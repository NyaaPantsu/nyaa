package api

// CreateUserAuthentication creates user authentication.
func CreateUserAuthentication(form *formStruct.LoginForm) (models.User, int, error) {
	username := form.Username
	pass := form.Password
	user, status, err := user.Exists(username, pass)
	return user, status, err
}