package errors

import (
	"fmt"
	"net/http"
)

var _ error = (*Err)(nil)

// Err represents an application error.
type Err struct {
	msg        string
	statusCode int
	appCode    string
}

// Error returns the error message as string.
func (e Err) Error() string {
	return e.msg
}

// StatusCode returns the HTTP status code associated with
// the error.
func (e Err) StatusCode() int {
	return e.statusCode
}

// AppCode returns an internal code associated with the error,
// so that it can be more easily recognised by clients of the API.
func (e Err) AppCode() string {
	return e.appCode
}

// List of application errors
var (
	ErrBadRequest = Err{
		msg:        "bad request",
		statusCode: http.StatusBadRequest, appCode: "BadRequest"}
	ErrCannotRenderResponse = Err{
		msg:        "error rendering response",
		statusCode: http.StatusUnprocessableEntity, appCode: "CannotRenderResponse"}
	ErrNoDefaultVariant = Err{
		msg:        "no default variant defined for flag",
		statusCode: http.StatusUnprocessableEntity, appCode: "NoDefaultVariant"}
	ErrInvalidFlag = Err{
		msg:        "invalid flag",
		statusCode: http.StatusUnprocessableEntity, appCode: "InvalidFlag"}
	ErrNoVariantToDistribute = Err{
		msg:        "no variants to distribute, please check the rule configuration",
		statusCode: http.StatusUnprocessableEntity, appCode: "NoVariantToDistribute"}
	ErrNotFound = Err{
		msg:        "not found",
		statusCode: http.StatusNotFound, appCode: "NotFound"}
	ErrNotImplemented = Err{
		msg:        "not implemented",
		statusCode: http.StatusNotImplemented, appCode: "NotImplemented"}
)

// NotFound returns an ErrNotFound error, for the given entity.
func NotFound(entity string) error {
	return fmt.Errorf("%s: %w", entity, ErrNotFound)
}

// BadRequest returns an ErrBadRequest error, with an additional message.
func BadRequest(message string) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, message)
}

// InvalidFlag returns an ErrInvalidFlag error, with an additional message.
func InvalidFlag(message string) error {
	return fmt.Errorf("%w: %s", ErrInvalidFlag, message)
}
