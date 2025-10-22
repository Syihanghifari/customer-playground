package usecase_customer

import (
	"context"
	"customer-playground/domain"
	"customer-playground/types"
	"time"

	"github.com/sirupsen/logrus"
)

type customerUseCase struct {
	customerRepository domain.CustomerRepository
	logger             *logrus.Logger
}

func (c customerUseCase) GetAll(ctx context.Context) ([]domain.Customer, error) {
	customers, err := c.customerRepository.GetAll(ctx)
	if err != nil {
		c.logger.Errorf("customerUseCase/GetAll :%v", err)
		return nil, err
	}
	return customers, nil
}

func (c customerUseCase) GetByCustomerNumber(customerNumber int, ctx context.Context) (domain.Customer, error) {
	customer, err := c.customerRepository.GetByCustomerNumber(customerNumber, ctx)
	if err != nil {
		c.logger.Errorf("customerUseCase/GetByCustomerNumber :%v", err)
		return domain.Customer{}, err
	}

	return customer, nil
}

func (c customerUseCase) Insert(customer *domain.Customer, ctx context.Context) (domain.Response, error) {
	now := time.Now()

	if !customer.CreatedAt.Valid {
		customer.CreatedAt = types.NullTime{Time: now, Valid: true}
	}
	if !customer.UpdatedAt.Valid {
		customer.UpdatedAt = types.NullTime{Time: now, Valid: true}
	}

	mesage, err := c.customerRepository.Insert(customer, ctx)
	if err != nil {
		c.logger.Errorf("customerUseCase/Insert :%v", err)
		return mesage, err
	}

	return mesage, nil
}

func (c customerUseCase) Update(newCustomer *domain.Customer, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	now := time.Now()
	currentCustomer, err := c.customerRepository.GetByCustomerNumber(newCustomer.CustomerNumber, ctx)
	if err != nil {
		message.Message = "this number is not found"
		message.StatusCode = 500
		return message, err
	}

	if newCustomer.Name == "" {
		newCustomer.Name = currentCustomer.Name
	}
	if newCustomer.Email == "" {
		newCustomer.Email = currentCustomer.Email
	}
	if newCustomer.Phone == "" {
		newCustomer.Phone = currentCustomer.Phone
	}
	if !newCustomer.BirthDate.Valid {
		newCustomer.BirthDate = currentCustomer.BirthDate
	}
	newCustomer.CreatedAt = currentCustomer.CreatedAt
	newCustomer.UpdatedAt = types.NullTime{Time: now, Valid: true}

	message, err = c.customerRepository.Update(newCustomer, ctx)
	if err != nil {
		c.logger.Errorf("customerUseCase/Update :%v", err)
		return message, err
	}
	return message, nil
}

func (c customerUseCase) DeleteByCustomerNumber(customerNumber int, ctx context.Context) (domain.Response, error) {
	message, err := c.customerRepository.DeleteByCustomerNumber(customerNumber, ctx)
	if err != nil {
		c.logger.Errorf("customerUseCase/DeleteByCustomerNumber :%v", err)
		return message, err
	}
	return message, nil
}

func NewCustomerUseCase(c domain.CustomerRepository, log *logrus.Logger) domain.CustomerUseCase {
	return &customerUseCase{
		customerRepository: c,
		logger:             log,
	}
}
