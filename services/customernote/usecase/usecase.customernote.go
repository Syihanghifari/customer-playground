package usecase_customernote

import (
	"context"
	"customer-playground/domain"
	"customer-playground/types"
	"time"

	"github.com/sirupsen/logrus"
)

type customerNoteUseCase struct {
	customerNoteRepository domain.CustomerNoteRepository
	logger                 *logrus.Logger
}

func (c customerNoteUseCase) GetAll(ctx context.Context) ([]domain.CustomerNote, error) {
	customerNotes, err := c.customerNoteRepository.GetAll(ctx)
	if err != nil {
		c.logger.Errorf("customerNoteUseCase/GetAll :%v", err)
		return nil, err
	}
	return customerNotes, nil
}
func (c customerNoteUseCase) GetByCustomerNumber(customerNumber int, ctx context.Context) ([]domain.CustomerNote, error) {
	customerNotes, err := c.customerNoteRepository.GetByCustomerNumber(customerNumber, ctx)
	if err != nil {
		c.logger.Errorf("customerNoteUseCase/GetByCustomerNumber :%v", err)
		return nil, err
	}
	return customerNotes, nil
}

func (c customerNoteUseCase) GetById(id int, ctx context.Context) (domain.CustomerNote, error) {
	customerNote, err := c.customerNoteRepository.GetById(id, ctx)
	if err != nil {
		c.logger.Errorf("customerNoteUseCase/GetById :%v", err)
		return domain.CustomerNote{}, err
	}
	return customerNote, nil
}

func (c customerNoteUseCase) Insert(customerNote *domain.CustomerNote, ctx context.Context) (domain.Response, error) {
	now := time.Now()
	if !customerNote.CreatedAt.Valid {
		customerNote.CreatedAt = types.NullTime{Time: now, Valid: true}
	}
	message, err := c.customerNoteRepository.Insert(customerNote, ctx)
	if err != nil {
		c.logger.Errorf("customerNoteUseCase/Insert :%v", err)
		return message, err
	}
	return message, nil
}

func (c customerNoteUseCase) Update(newCustomerNote *domain.CustomerNote, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	currentCustomerNote, err := c.customerNoteRepository.GetById(newCustomerNote.CustomerNumber, ctx)
	if err != nil {
		message.Message = "this id is not found"
		message.StatusCode = 500
		return message, err
	}

	if newCustomerNote.CustomerNumber == 0 {
		newCustomerNote.CustomerNumber = currentCustomerNote.CustomerNumber
	}
	if newCustomerNote.Note == "" {
		newCustomerNote.Note = currentCustomerNote.Note
	}
	if !newCustomerNote.CreatedAt.Valid {
		newCustomerNote.CreatedAt = currentCustomerNote.CreatedAt
	}
	message, err = c.customerNoteRepository.Update(newCustomerNote, ctx)
	if err != nil {
		c.logger.Errorf("customerNoteUseCase/Update :%v", err)
		return message, err
	}
	return message, nil
}

func (c customerNoteUseCase) DeleteById(id int, ctx context.Context) (domain.Response, error) {
	message, err := c.customerNoteRepository.DeleteById(id, ctx)
	if err != nil {
		c.logger.Errorf("customerNoteUseCase/DeleteById/DelCustomerNumber :%v", err)
		return message, err
	}
	return message, nil
}

func NewCustomerNoteUseCase(c domain.CustomerNoteRepository, log *logrus.Logger) domain.CustomerNoteUseCase {
	return &customerNoteUseCase{
		customerNoteRepository: c,
		logger:                 log,
	}
}
