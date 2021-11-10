package swagger

type UserRequest struct {
	Email    string `json:"email" validate:"required" example:"test@email.test"`
	Password string `json:"password" validate:"required" example:"p4sSWorD"`
}
