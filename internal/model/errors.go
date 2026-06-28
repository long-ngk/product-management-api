package model

// ValidationError represents a validation failure in the Service layer.
// Handler maps this to HTTP 400 Bad Request.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NotFoundError represents a resource not found error in the Service layer.
// Handler maps this to HTTP 404 Not Found.
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
