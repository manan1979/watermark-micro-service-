package transport

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/manan1979/watermark-service/internal/util"
	"github.com/manan1979/watermark-service/pkg/auth/endpoint"
)

func NewHTTPHandler(ep endpoint.AuthSet) http.Handler {
	m := http.NewServeMux()

	m.Handle("/user/access", httptransport.NewServer(
		ep.GetUserAccessEndpoint,
		decodeHTTPGetUserAccessRequest,
		encodeResponse,
	))

	m.Handle("/authenticate", httptransport.NewServer(
		ep.AuthenticateEndpoint,
		decodeHTTPGetUserAuthenticateRequest,
		encodeResponse,
	))

	m.Handle("/healthz", httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))

	return m
}

func decodeHTTPGetUserAccessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.GetUserAccessRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPGetUserAuthenticateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.AuthenticateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.ServiceStatusRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case util.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)

	case util.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
