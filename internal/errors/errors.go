package errors

import (
	"fmt"
	"net/http"
)

var _ error = (*Err)(nil)

type Err struct {
	msg        string
	statusCode int
	appCode    string
}

func (e Err) Error() string {
	return e.msg
}

func (e Err) StatusCode() int {
	return e.statusCode
}

func (e Err) AppCode() string {
	return e.appCode
}

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

func NotFound(entity string) error {
	return fmt.Errorf("%s: %w", entity, ErrNotFound)
}

func BadRequest(message string) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, message)
}
