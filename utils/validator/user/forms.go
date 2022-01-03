package userValidator

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username           string `validate:"required,min=3,max=20" form:"username" json:"username"`
	Email              string `form:"email" json:"email"`
	Password           string `validate:"required,min=6,max=72,eqfield=ConfirmPassword" form:"password" json:"password"`
	ConfirmPassword    string `validate:"required" omit:"true" form:"password_confirmation" json:"password_confirmation"` // Omit when binding to user model since user model doesn't have those field
	CaptchaID          string `validate:"required" omit:"true" form:"captchaID" json:"captchaID"`
	TermsAndConditions string `validate:"eq=1" omit:"true" form:"t_and_c" json:"t_and_c"`
}

// LoginForm is used when a user logs in.
type LoginForm struct {
	Username   string `validate:"required" json:"username" form:"username"`
	Password   string `validate:"required" json:"password" form:"password"`
	RedirectTo string `validate:"-" form:"redirectTo" json:"-"`
	RememberMe string `validate:"-" form:"remember_me" json:"-"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Username        string `validate:"required" form:"username" json:"username" needed:"true" len_min:"3" len_max:"20"`
	Email           string `json:"email" form:"email"`
	Language        string `validate:"default=en-us" form:"language" json:"language"`
	CurrentPassword string `validate:"omitempty,min=6,max=72" form:"current_password" json:"current_password" omit:"true"`
	Password        string `validate:"omitempty,min=6,max=72" form:"password" json:"password" len_min:"6" len_max:"72" equalInput:"ConfirmPassword"`
	ConfirmPassword string `validate:"omitempty" form:"password_confirmation" json:"password_confirmation" omit:"true"`
	Status          int    `validate:"default=0" form:"status" json:"status"`
	Theme           string `form:"theme" json:"theme"`
	AnidexAPIToken  string `validate:"-" form:"anidex_api" json:"anidex_api"`
	NyaasiAPIToken  string `validate:"-" form:"nyaasi_api" json:"nyaasi_api"`
	TokyoTAPIToken  string `validate:"-" form:"tokyot_api" json:"tokyot_api"`
}

// UserSettingsForm is used when updating a user.
type UserSettingsForm struct {
	NewTorrent        bool `validate:"-" json:"new_torrent" form:"new_torrent"`
	NewTorrentEmail   bool `validate:"-" json:"new_torrent_email" form:"new_torrent_email"`
	NewComment        bool `validate:"-" json:"new_comment" form:"new_comment"`
	NewCommentEmail   bool `validate:"-" json:"new_comment_email" form:"new_comment_email"`
	NewResponses      bool `validate:"-" json:"new_responses" form:"new_responses"`
	NewResponsesEmail bool `validate:"-" json:"new_responses_email" form:"new_responses_email"`
	NewFollower       bool `validate:"-" json:"new_follower" form:"new_follower"`
	NewFollowerEmail  bool `validate:"-" json:"new_follower_email" form:"new_follower_email"`
	Followed          bool `validate:"-" json:"followed" form:"followed"`
	FollowedEmail     bool `validate:"-" json:"followed_email" form:"followed_email"`
}

// PasswordForm is used when updating a user password.
type PasswordForm struct {
	CurrentPassword string `validate:"required" form:"current_password"`
	Password        string `validate:"required" form:"password"`
}

// SendPasswordResetForm is used when sending a password reset token.
type SendPasswordResetForm struct {
	Email string `validate:"required" form:"email"`
}

// PasswordResetForm is used when reseting a password.
type PasswordResetForm struct {
	PasswordResetToken string `validate:"required" form:"token"`
	Password           string `validate:"required" form:"password"`
}
