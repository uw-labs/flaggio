package service

import (
	"crypto/sha1" // nolint // only used for hashing requests
	"encoding/hex"
	"net/http"
	"sort"

	"github.com/victorkt/clientip"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/vmihailenco/msgpack/v4"
)

// EvaluationRequest is the evaluation request object
type EvaluationRequest struct {
	UserID      string              `json:"userId"`
	UserContext flaggio.UserContext `json:"context"`
	Debug       *bool               `json:"debug,omitempty"`
}

// Bind adds additional data to the EvaluationRequest.
// Some special fields are added to the user context:
// * $userId is the user ID provided in the request
// * $ip is the network address that originated the request
func (er EvaluationRequest) Bind(r *http.Request) error {
	// enrich user context
	er.UserContext["$userId"] = er.UserID
	er.UserContext["$ip"] = clientip.FromContext(r.Context()).String()
	return nil
}

// IsDebug returns whether this is a debug request or not
func (er EvaluationRequest) IsDebug() bool {
	return er.Debug != nil && *er.Debug
}

// Hash returns a hash string representation of EvaluationRequest
// This function will return the same hash regardless of the order
// the user context comes in.
func (er EvaluationRequest) Hash() (string, error) {
	// sort user context keys
	var contextKeys []string
	for key := range er.UserContext {
		contextKeys = append(contextKeys, key)
	}
	sort.Strings(contextKeys)

	// create 2d slice with sorted keys from user context
	ordered := make([]interface{}, len(contextKeys))
	for idx, key := range contextKeys {
		ordered[idx] = []interface{}{key, er.UserContext[key]}
	}

	// marshal ordered slice and hash it
	bytes, err := msgpack.Marshal(ordered)
	if err != nil {
		return "", err
	}
	h := sha1.New() // nolint // we don't care about security for this
	if _, err := h.Write(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// EvaluationResponse is the evaluation response object
type EvaluationResponse struct {
	Evaluation  *flaggio.Evaluation  `json:"evaluation"`
	UserContext *flaggio.UserContext `json:"context,omitempty"`
}

// Render can enrich the EvaluationResponse object before being returned to the
// user. Currently it does nothing, but is needed to satisfy the
// chi.Renderer interface.
func (e *EvaluationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// EvaluationsResponse is the evaluation response object
type EvaluationsResponse struct {
	Evaluations flaggio.EvaluationList `json:"evaluations"`
	UserContext *flaggio.UserContext   `json:"context,omitempty"`
}

// Render can enrich the EvaluationsResponse object before being returned to the
// user. Currently it does nothing, but is needed to satisfy the
// chi.Renderer interface.
func (e *EvaluationsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
