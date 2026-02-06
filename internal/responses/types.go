package responses

type Success[T any] struct {
	Code    int     `json:"code"`
	Message *string `json:"message,omitempty"`
	Data    *T      `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
