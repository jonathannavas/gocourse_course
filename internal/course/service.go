package course

import (
	"context"
	"log"
	"time"

	"github.com/jonathannavas/gocourse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name *string, startDate, endDate *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	if startDateParsed.After(endDateParsed) {
		s.log.Println(errDateValidation)
		return nil, errDateValidation
	}

	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	courses, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {
	course, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s service) Update(ctx context.Context, id string, name *string, startDate, endDate *string) error {

	var startDateParsed, endDateParsed *time.Time

	course, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Println(err)
		return err
	}

	if startDate != nil {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Println(err)
			return err
		}

		if date.After(course.EndDate) {
			s.log.Println(errDateValidation)
			return errDateValidation
		}

		startDateParsed = &date
		*startDateParsed = startDateParsed.Add(time.Hour * 24)

	}

	if endDate != nil {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		if course.StartDate.After(date) {
			s.log.Println(errDateValidation)
			return errDateValidation
		}
		endDateParsed = &date
		*endDateParsed = endDateParsed.Add(time.Hour * 24)
	}

	if err := s.repo.Update(ctx, id, name, startDateParsed, endDateParsed); err != nil {
		return err
	}

	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
