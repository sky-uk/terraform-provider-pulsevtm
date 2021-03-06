package api

// VTMError : Generic error object for Brocade vTM
type VTMError struct {
	ErrorID   string                 `json:"error_id"`
	ErrorText string                 `json:"error_text,omitempty"`
	ErrorInfo map[string]interface{} `json:"error_info,omitempty"`
}
