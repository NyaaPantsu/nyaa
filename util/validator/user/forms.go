package userValidator

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username           string `form:"username" needed:"true" len_min:"3" len_max:"20"`
	Email              string `form:"email"`
	Password           string `form:"password" needed:"true" len_min:"6" len_max:"72" equalInput:"ConfirmPassword"`
	ConfirmPassword    string `form:"password_confirmation" hum_name:"Password Confirmation" omit:"true" needed:"true"`
	CaptchaID          string `form:"captchaID" omit:"true" needed:"true"`
	TermsAndConditions bool   `form:"t_and_c" omit:"true" needed:"true" equal:"true" hum_name:"Terms and Conditions"`
}

// LoginForm is used when a user logs in.
type LoginForm struct {
	Username string `form:"username" needed:"true" json:"username"`
	Password string `form:"password" needed:"true" json:"password"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Username        string `form:"username" json:"username" needed:"true" len_min:"3" len_max:"20"`
	Email           string `form:"email" json:"email"`
	Language        string `form:"language" json:"language" default:"en-us"`
	CurrentPassword string `form:"current_password" json:"current_password" len_min:"6" len_max:"72" omit:"true"`
	Password        string `form:"password" json:"password" len_min:"6" len_max:"72" equalInput:"ConfirmPassword"`
	ConfirmPassword string `form:"password_confirmation" json:"password_confirmation" hum_name:"Password Confirmation" omit:"true"`
	Status          int    `form:"status" json:"status" default:"0"`
	Theme           string `form:"theme" json:"theme"`
}

// UserSettingsForm is used when updating a user.
type UserSettingsForm struct {
	NewTorrent        bool `form:"new_torrent" json:"new_torrent" default:"true"`
	NewTorrentEmail   bool `form:"new_torrent_email" json:"new_torrent_email" default:"true"`
	NewComment        bool `form:"new_comment" json:"new_comment" default:"true"`
	NewCommentEmail   bool `form:"new_comment_email" json:"new_comment_email" default:"false"`
	NewResponses      bool `form:"new_responses" json:"new_responses" default:"true"`
	NewResponsesEmail bool `form:"new_responses_email" json:"new_responses_email" default:"false"`
	NewFollower       bool `form:"new_follower" json:"new_follower" default:"true"`
	NewFollowerEmail  bool `form:"new_follower_email" json:"new_follower_email" default:"true"`
	Followed          bool `form:"followed" json:"followed" default:"false"`
	FollowedEmail     bool `form:"followed_email" json:"followed_email" default:"false"`
}

// PasswordForm is used when updating a user password.
type PasswordForm struct {
	CurrentPassword string `form:"currentPassword" needed:"true"`
	Password        string `form:"newPassword" needed:"true"`
}

// SendPasswordResetForm is used when sending a password reset token.
type SendPasswordResetForm struct {
	Email string `form:"email" needed:"true"`
}

// PasswordResetForm is used when reseting a password.
type PasswordResetForm struct {
	PasswordResetToken string `form:"token" needed:"true"`
	Password           string `form:"newPassword" needed:"true"`
}
