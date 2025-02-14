package responses

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Standard success/error response functions
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func NewErrorResponse(err string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
	}
}
