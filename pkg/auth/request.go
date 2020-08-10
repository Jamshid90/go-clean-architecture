package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignupRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required"`
	Gender          string `json:"gender" validate:"required,eq=male|eq=female"`
	FirstName       string `json:"first_name" validate:"required,min=2,max=50"`
	LastName        string `json:"last_name" validate:"required,min=2,max=50"`
	BirthDate       string `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,eqfield=Password"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
