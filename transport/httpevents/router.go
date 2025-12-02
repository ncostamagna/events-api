package httpevents

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ncostamagna/events-api/pkg/httputil"
	"github.com/ncostamagna/go-http-utils/response"
)

type ctxKey string

const (
	ctxParam  ctxKey = "params"
	ctxHeader ctxKey = "header"
	ctxQuery  ctxKey = "query"
)

func NewHTTPServer(endpoints Endpoints) http.Handler {

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(ginDecode())

	r.GET("/events", gin.WrapH(httputil.NewServer(endpoints.GetAll, decodeGetAllHandler, encodeResponse, encodeError)))
	r.POST("/events", gin.WrapH(httputil.NewServer(endpoints.Store, decodeStoreHandler, encodeResponse, encodeError)))
	r.GET("/events/:id", gin.WrapH(httputil.NewServer(endpoints.Get, decodeGetHandler, encodeResponse, encodeError)))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func ginDecode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ctxParam, c.Params)
		ctx = context.WithValue(ctx, ctxHeader, c.Request.Header)
		ctx = context.WithValue(ctx, ctxQuery, c.Request.URL.Query())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func decodeGetHandler(ctx context.Context, _ *http.Request) (interface{}, error) {
	params := ctx.Value(ctxParam).(gin.Params)

	return GetReq{
		ID: params.ByName("id"),
	}, nil
}

func decodeGetAllHandler(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeStoreHandler(_ context.Context, r *http.Request) (interface{}, error) {
	var req StoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
