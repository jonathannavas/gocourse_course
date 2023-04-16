package course

import (
	"context"
	"errors"

	"github.com/jonathannavas/go_lib_response/response"
	"github.com/jonathannavas/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	GetReq struct {
		ID string
	}

	DeleteReq struct {
		ID string
	}

	GetAllReq struct {
		Name     string
		LastName string
		Limit    int
		Page     int
	}

	UpdateRequest struct {
		ID        string
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	// Response struct {
	// 	Status int         `json:"status"`
	// 	Data   interface{} `json:"data,omitempty"`
	// 	Error  string      `json:"error,omitempty"`
	// 	Meta   *meta.Meta  `json:"meta,omitempty"`
	// }

	Config struct {
		LimitPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		courseBody := request.(CreateRequest)

		if courseBody.Name == "" {
			return nil, response.BadRequest(errNameRequired.Error())
		}

		if courseBody.StartDate == "" {
			return nil, response.BadRequest(errStartDateRequired.Error())
		}

		if courseBody.EndDate == "" {
			return nil, response.BadRequest(errEndDateRequired.Error())
		}

		course, err := s.Create(ctx, courseBody.Name, courseBody.StartDate, courseBody.EndDate)

		if err != nil {
			if err == errDateValidation {
				return nil, response.BadRequest(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil, 200), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		course, err := s.Get(ctx, req.ID)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", course, nil, 201), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			Name: req.Name,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimitPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", courses, meta, 201), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(errNameRequired.Error())
		}
		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest(errStartDateRequired.Error())
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(errStartDateRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.Name, req.StartDate, req.EndDate)

		if err == errDateValidation {
			return nil, response.BadRequest(err.Error())
		}

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil, 201), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil, 200), nil
	}
}
