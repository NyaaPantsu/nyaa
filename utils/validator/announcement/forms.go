package announcementValidator

// CreateForm is a struct to validate the anouncement before adding to db
type CreateForm struct {
	ID      uint   `validate:"-"`
	Message string `validate:"required,min=5" form:"message"`
	Delay   int    `validate:"omitempty,min=1" form:"delay"`
}
