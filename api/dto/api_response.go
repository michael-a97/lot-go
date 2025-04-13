package dto

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Paging  any    `json:"paging,omitempty"`
	Error   any    `json:"error"`
}
