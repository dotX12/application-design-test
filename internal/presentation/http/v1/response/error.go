package response

type ErrorDetailComponentResponse struct {
	Field []string `json:"field"   example:"source"`                                          // Field that caused the error (if any)
	Error string   `json:"message" example:"source: Field valid values are [foo bar] (enum)"` // Error message
} // @name ErrorDetailComponentResponse

type ErrorDetailsResponse struct {
	Message string                         `json:"message" example:"Validation failed for the request parameters"` // Detailed message of the error
	Status  string                         `json:"status"  example:"Unprocessable Entity"`                         // Small message of the error
	Slug    string                         `json:"slug"    example:"400_malformed_request"`                        // Slug of the error
	Details []ErrorDetailComponentResponse `json:"details"`                                                        // Details of the error
} // @name ErrorDetailsResponse

type ErrorResponse struct {
	Error ErrorDetailsResponse `json:"error"` // Error response
} // @name ErrorResponse
