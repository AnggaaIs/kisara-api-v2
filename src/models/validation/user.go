package validation

type UserUpdateBody struct {
	Nickname string `json:"nickname" body:"nickname" validate:"omitempty,max=15"`
}
