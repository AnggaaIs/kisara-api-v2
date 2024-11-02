package validation

type AuthGoogleCallbackBody struct {
	Code  string `json:"code" body:"code" validate:"required"`
	State string `json:"state" body:"state" validate:"required"`
}
