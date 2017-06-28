package users

// UpdateUserCore updates a user. (Applying the modifed data of user).
func Update(user *model.User) (int, error) {
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

// UpdateRawUser : Function to update a user without updating his associations model
func UpdateRawUser(user *model.User) (int, error) {
	user.UpdatedAt = time.Now()
	err := db.ORM.Model(&user).UpdateColumn(&user).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}