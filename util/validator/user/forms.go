package userValidator

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username           string `validate:"required,min=3,max=20"`
	Email              string
	Password           string `validate:"required,min=6,max=72,eqfield=ConfirmPassword"`
	ConfirmPassword    string `validate:"required" omit:"true"` // Omit when binding to user model since user model doesn't have those field
	CaptchaID          string `validate:"required" omit:"true"`
	TermsAndConditions string `validate:"eq=true" omit:"true"`
}

// LoginForm is used when a user logs in.
type LoginForm struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Username        string `validate:"required" json:"username" needed:"true" len_min:"3" len_max:"20"`
	Email           string `json:"email"`
	Language        string `validate:"default=en-us" json:"language"`
	CurrentPassword string `validate:"required,min=6,max=72" json:"current_password" omit:"true"`
	Password        string `validate:"required,min=6,max=72" json:"password" len_min:"6" len_max:"72" equalInput:"ConfirmPassword"`
	ConfirmPassword string `validate:"required" json:"password_confirmation" omit:"true"`
	Status          int    `validate:"default=0" json:"status"`
	Theme           string `json:"theme"`
}

// UserSettingsForm is used when updating a user.
type UserSettingsForm struct {
	NewTorrent        bool `validate:"default=true" json:"new_torrent"`
	NewTorrentEmail   bool `validate:"default=true" json:"new_torrent_email"`
	NewComment        bool `validate:"default=true" json:"new_comment"`
	NewCommentEmail   bool `validate:"default=false" json:"new_comment_email"`
	NewResponses      bool `validate:"default=true" json:"new_responses"`
	NewResponsesEmail bool `validate:"default=false" json:"new_responses_email"`
	NewFollower       bool `validate:"default=true" json:"new_follower"`
	NewFollowerEmail  bool `validate:"default=true" json:"new_follower_email"`
	Followed          bool `validate:"default=false" json:"followed"`
	FollowedEmail     bool `validate:"default=false" json:"followed_email"`
}

// PasswordForm is used when updating a user password.
type PasswordForm struct {
	CurrentPassword string `validate:"required"`
	Password        string `validate:"required"`
}

// SendPasswordResetForm is used when sending a password reset token.
type SendPasswordResetForm struct {
	Email string `validate:"required"`
}

// PasswordResetForm is used when reseting a password.
type PasswordResetForm struct {
	PasswordResetToken string `validate:"required"`
	Password           string `validate:"required"`
}
