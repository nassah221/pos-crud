package dto

type GenericResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}
