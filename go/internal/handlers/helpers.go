package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// Common error codes matching TypeSpec definition
const (
	ErrorCodeBadRequest             = generated.BADREQUEST
	ErrorCodeUnauthorized           = generated.UNAUTHORIZED
	ErrorCodeForbidden              = generated.FORBIDDEN
	ErrorCodeNotFound               = generated.NOTFOUND
	ErrorCodeConflict               = generated.CONFLICT
	ErrorCodeValidationError        = generated.VALIDATIONERROR
	ErrorCodeInsufficientStock      = generated.INSUFFICIENTSTOCK
	ErrorCodeInvalidStateTransition = generated.INVALIDSTATETRANSITION
	ErrorCodeInternalError          = generated.INTERNALERROR
	ErrorCodeServiceUnavailable     = generated.SERVICEUNAVAILABLE
)

// errorResponse sends a standardized error response
func errorResponse(w http.ResponseWriter, statusCode int, code generated.ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := generated.ErrorResponse{
		Error: struct {
			Code    generated.ErrorCode `json:"code"`
			Details *interface{}        `json:"details,omitempty"`
			Message string              `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}
	json.NewEncoder(w).Encode(response)
}
