package domain

import (
	"context"
	"customer-playground/types"
)

type CustomerNote struct {
	ID             int            `json:"id"`
	CustomerNumber int            `json:"customer_number"`
	Note           string         `json:"note"`
	CreatedAt      types.NullTime `json:"created_at,omitempty" swaggertype:"string" example:"1995-06-12T00:00:00Z"`
}

type (
	CustomerNoteUseCase interface {
		GetAll(ctx context.Context) ([]CustomerNote, error)
		GetByCustomerNumber(customerNumber int, ctx context.Context) ([]CustomerNote, error)
		GetById(id int, ctx context.Context) (CustomerNote, error)
		Insert(customerNote *CustomerNote, ctx context.Context) (Response, error)
		Update(customerNote *CustomerNote, ctx context.Context) (Response, error)
		DeleteById(id int, ctx context.Context) (Response, error)
	}
	CustomerNoteRepository interface {
		GetAll(ctx context.Context) ([]CustomerNote, error)
		GetByCustomerNumber(customerNumber int, ctx context.Context) ([]CustomerNote, error)
		GetById(id int, ctx context.Context) (CustomerNote, error)
		Insert(customerNote *CustomerNote, ctx context.Context) (Response, error)
		Update(customerNote *CustomerNote, ctx context.Context) (Response, error)
		DeleteById(id int, ctx context.Context) (Response, error)
	}
)
