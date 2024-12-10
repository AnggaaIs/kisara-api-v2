package validation

type MessageBodyPost struct {
	MessageContent string `json:"message_content" body:"message_content" validate:"required"`
}

type MessageBodyGet struct {
	SortBy string `query:"sort_by" validate:"omitempty,oneof=asc desc" default:"desc"`
	Page   int    `query:"page" validate:"omitempty,gte=1" default:"1"`
	Limit  int    `query:"limit" validate:"omitempty,gte=1" default:"10"`
}
