package validation

type MessageBodyPost struct {
	MessageContent string `json:"message_content" body:"message_content" validate:"required"`
}
