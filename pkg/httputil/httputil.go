package httputil

import (
	"context"
	"net/http"
)

type Endpoint func(ctx context.Context, request interface{}) (interface{}, error)
type DecodeRequestFunc func(context.Context, *http.Request) (interface{}, error)
type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error
type EncodeErrorFunc func(context.Context, error, http.ResponseWriter)

func NewServer(
	endpoint Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	encodeError EncodeErrorFunc,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request, err := dec(ctx, r)
		if err != nil {
			encodeError(ctx, err, w)
			return
		}

		response, err := endpoint(ctx, request)
		if err != nil {
			encodeError(ctx, err, w)
			return
		}

		if err := enc(ctx, w, response); err != nil {
			encodeError(ctx, err, w)
			return
		}
	})
}

func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == http.MethodOptions {
			return
		}

		h.ServeHTTP(w, r)
	})
}
