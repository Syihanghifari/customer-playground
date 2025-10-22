package repository_customer

import (
	"context"
	"customer-playground/domain"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type customerRepository struct {
	dbPool *sql.DB
	logger *logrus.Logger
}

func (c customerRepository) GetAll(ctx context.Context) ([]domain.Customer, error) {

	stmt, err := c.dbPool.Prepare(`
		SELECT
			customer_number,
			name,
			email,
			phone,
			birth_date,
			created_at,
			updated_at
		FROM customer
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var customer domain.Customer
		err := rows.Scan(
			&customer.CustomerNumber,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.BirthDate,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (c customerRepository) GetByCustomerNumber(customerNumber int, ctx context.Context) (domain.Customer, error) {
	stmt, err := c.dbPool.Prepare(`
		SELECT
			customer_number,
			name,
			email,
			phone,
			birth_date,
			created_at,
			updated_at
		FROM customer
		WHERE customer_number = $1
		ORDER BY name
	`)
	if err != nil {
		return domain.Customer{}, err
	}
	var customer domain.Customer
	err = stmt.QueryRowContext(ctx, customerNumber).Scan(
		&customer.CustomerNumber,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.BirthDate,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}

func (c customerRepository) Insert(customer *domain.Customer, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	stmt, err := c.dbPool.PrepareContext(ctx, `
		INSERT INTO customer(
			customer_number,
			name,
			email,
			phone,
			birth_date,
			created_at,
			updated_at) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	)
	`)
	if err != nil {
		message.Message = "Failed to Insert Customer"
		message.StatusCode = 500
		c.logger.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		customer.CustomerNumber,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.BirthDate,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	if err != nil {
		message.Message = "Failed to Insert Customer Note"
		message.StatusCode = 500
		c.logger.Errorf("failed to execute statement: %v", err)
	}
	message.Message = fmt.Sprintf("Succes Insert Customer with number %d", customer.CustomerNumber)
	message.StatusCode = 200
	return message, nil
}

func (c customerRepository) Update(customer *domain.Customer, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	stmt, err := c.dbPool.PrepareContext(ctx, `
		UPDATE customer SET
			name = $2,
			email = $3,
			phone = $4,
			birth_date = $5, 
			created_at = $6,
			updated_at = $7
		WHERE 
			customer_number = $1
	`)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Update number %d", customer.CustomerNumber)
		message.StatusCode = 500
		c.logger.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		customer.CustomerNumber,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.BirthDate,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Update id %d", customer.CustomerNumber)
		message.StatusCode = 500
		c.logger.Errorf("failed to execute statement: %v", err)
	}
	message.Message = "Succes Update"
	message.StatusCode = 200
	return message, nil
}

func (c customerRepository) DeleteByCustomerNumber(customerNumber int, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	stmt, err := c.dbPool.Prepare(`
		DELETE
		FROM customer
		WHERE customer_number = $1
	`)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Delete number %d", customerNumber)
		message.StatusCode = 500
		return message, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, customerNumber)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Delete number %d", customerNumber)
		message.StatusCode = 500
		return message, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		message.Message = fmt.Sprintf("Could not determine rows affected for number %d", customerNumber)
		message.StatusCode = 500
		c.logger.Errorf("failed to get rows affected: %v", err)
		return message, err
	}
	if rowsAffected == 0 {
		message.Message = fmt.Sprintf("No customer note found with number %d", customerNumber)
		message.StatusCode = 500
		return message, sql.ErrNoRows
	}
	message.StatusCode = 200
	message.Message = "Succes Delete!!"
	return message, nil
}

func NewCustomerRepository(db *sql.DB, log *logrus.Logger) domain.CustomerRepository {
	return &customerRepository{
		dbPool: db,
		logger: log,
	}
}
