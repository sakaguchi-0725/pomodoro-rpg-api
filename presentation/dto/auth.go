package dto

type SignInRequest struct {
	Email    string `json:"email"`
	Passowrd string `json:"password"`
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmSignUpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ChangePasswordRequest struct {
	PreviousPass string `json:"previousPassword"`
	ProposedPass string `json:"proposedPassword"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ConfirmForgotPasswordRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
