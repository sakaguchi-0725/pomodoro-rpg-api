package dto

type AccountResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type UpdateAccountRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
