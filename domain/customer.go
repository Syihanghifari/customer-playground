package domain

import (
	"context"
	"customer-playground/types"
)

type Customer struct {
	CustomerNumber int            `json:"customer_number"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Phone          string         `json:"phone,omitempty"`
	BirthDate      types.NullTime `json:"birth_date,omitempty" swaggertype:"string" example:"1995-06-12T00:00:00Z"`
	CreatedAt      types.NullTime `json:"created_at,omitempty" swaggertype:"string" example:"1995-06-12T00:00:00Z"`
	UpdatedAt      types.NullTime `json:"updated_at,omitempty" swaggertype:"string" example:"1995-06-12T00:00:00Z"`
}

type (
	CustomerUseCase interface {
		GetAll(ctx context.Context) ([]Customer, error)
		GetByCustomerNumber(customerNumber int, ctx context.Context) (Customer, error)
		Insert(customer *Customer, ctx context.Context) (Response, error)
		Update(customer *Customer, ctx context.Context) (Response, error)
		DeleteByCustomerNumber(customerNumber int, ctx context.Context) (Response, error)
	}

	CustomerRepository interface {
		GetAll(ctx context.Context) ([]Customer, error)
		GetByCustomerNumber(customerNumber int, ctx context.Context) (Customer, error)
		Insert(customer *Customer, ctx context.Context) (Response, error)
		Update(customer *Customer, ctx context.Context) (Response, error)
		DeleteByCustomerNumber(customerNumber int, ctx context.Context) (Response, error)
	}
)
