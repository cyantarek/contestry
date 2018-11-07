package auth

type RegistrationPayload struct {
	FirstName       string `form:"f_name" binding:"required"`
	LastName        string `form:"l_name" binding:"required"`
	Username        string `form:"username" binding:"required"`
	ID              string `form:"id" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" binding:"required"`
	UserType        string `form:"user_type" binding:"required"`
}

type LoginPayload struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
