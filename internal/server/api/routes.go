package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	internalerrors "github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/service"
)

// POST /evaluate/{id}
// Evaluates a given flag for the user
func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "POST /evaluate/{id}")
	defer span.Finish()

	flagKey := chi.URLParam(r, "key")
	er := &service.EvaluationRequest{
		UserContext: make(flaggio.UserContext),
	}
	defer r.Body.Close()

	// unmarshal JSON request
	if err := render.Bind(r, er); err != nil {
		badRequest := internalerrors.BadRequest(err.Error())
		_ = render.Render(w, r, formatErr(badRequest))
		return
	}

	// evaluate flag
	eval, err := s.flagsService.Evaluate(ctx, flagKey, er)
	if err != nil {
		s.logger.WithError(err).WithField("req_id", middleware.GetReqID(ctx)).
			Error("failed to evaluate flag")
		_ = render.Render(w, r, formatErr(err))
		return
	}

	// render response
	if err = render.Render(w, r, eval); err != nil {
		cannotRender := fmt.Errorf("%w: %s", internalerrors.ErrCannotRenderResponse, err)
		_ = render.Render(w, r, formatErr(cannotRender))
		return
	}
}

// POST /evaluate
// Evaluates all flags for the user
func (s *Server) handleEvaluateAll(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "POST /evaluate")
	defer span.Finish()

	er := &service.EvaluationRequest{
		UserContext: make(flaggio.UserContext),
	}
	defer r.Body.Close()

	// unmarshal JSON request
	if err := render.Bind(r, er); err != nil {
		badRequest := internalerrors.BadRequest(err.Error())
		_ = render.Render(w, r, formatErr(badRequest))
		return
	}

	// evaluate flags
	eval, err := s.flagsService.EvaluateAll(ctx, er)
	if err != nil {
		s.logger.WithError(err).WithField("req_id", middleware.GetReqID(ctx)).
			Error("failed to evaluate all")
		_ = render.Render(w, r, formatErr(err))
		return
	}

	// render response
	if err = render.Render(w, r, eval); err != nil {
		cannotRender := fmt.Errorf("%w: %s", internalerrors.ErrCannotRenderResponse, err)
		_ = render.Render(w, r, formatErr(cannotRender))
		return
	}
}

type errResponse struct {
	Err        error  `json:"-"`               // low-level runtime error
	StatusCode int    `json:"-"`               // http response status code
	StatusText string `json:"status"`          // user-level status message
	AppCode    string `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func formatErr(err error) render.Renderer {
	res := &errResponse{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		StatusText: "error processing request",
		ErrorText:  err.Error(),
		AppCode:    "InternalServerError",
	}
	var e internalerrors.Err
	if errors.As(err, &e) {
		res.StatusCode = e.StatusCode()
		res.AppCode = e.AppCode()
	}
	return res
}
