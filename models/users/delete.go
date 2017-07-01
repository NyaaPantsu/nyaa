package users

// DeleteUser deletes a user.
func DeleteUser(currentUser *models.User, id string) (int, error) {
	var user models.User

	if db.ORM.First(&user, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("user_not_found")
	}
	if user.ID == 0 {
		return http.StatusInternalServerError, errors.New("permission_delete_error")
	}
	err := db.ORM.Delete(&user).Error
	if err != nil {
		return http.StatusInternalServerError, errors.New("user_not_deleted")
	}
	if user.CurrentUserIdentical(currentUser, user.ID) {
		ClearCookie(c)
	}

	return http.StatusOK, nil
}