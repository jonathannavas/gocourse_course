package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/jonathannavas/go_lib_response/response"
	"github.com/jonathannavas/gocourse_course/internal/course"
)

func NewUserHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {
	r := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	r.POST("/courses", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse, encodeResponse,
		opts...,
	)))

	r.GET("/courses", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourse, encodeResponse,
		opts...,
	)))

	r.GET("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse, encodeResponse,
		opts...,
	)))

	r.PATCH("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse, encodeResponse,
		opts...,
	)))

	r.DELETE("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse, encodeResponse,
		opts...,
	)))

	return r
}

func ginDecode(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "params", c.Params)
	c.Request = c.Request.WithContext(ctx)
}

func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var req course.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

func decodeGetCourse(ctx context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}
	params := ctx.Value("params").(gin.Params)
	req := course.GetReq{
		ID: params.ByName("id"),
	}
	return req, nil
}

func decodeGetAllCourse(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}
	v := r.URL.Query()
	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("limit"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

func decodeUpdateCourse(ctx context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var req course.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}
	params := ctx.Value("params").(gin.Params)
	req.ID = params.ByName("id")
	return req, nil
}

func decodeDeleteCourse(ctx context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	params := ctx.Value("params").(gin.Params)
	req := course.DeleteReq{
		ID: params.ByName("id"),
	}
	return req, nil
}

func authorization(token string) error {
	if token != os.Getenv("TOKEN") {
		return errors.New("invalid token")
	}
	return nil
}
