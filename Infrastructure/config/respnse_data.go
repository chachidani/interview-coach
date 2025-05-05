package config

type ResponseData struct {
	Error           bool        `json:"error"`
	SuccessResponse bool        `json:"successResponse"`
	ErrorMessage    string      `json:"errorMessage"`
	SuccessMessage  *string     `json:"successMessage"`
	Data            interface{} `json:"data,omitempty"`
}
