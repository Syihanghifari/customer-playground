package repository_customernote

import (
	"context"
	"customer-playground/domain"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type customerNoteRepository struct {
	dbPool *sql.DB
	logger *logrus.Logger
}

func (c customerNoteRepository) GetAll(ctx context.Context) ([]domain.CustomerNote, error) {
	stmt, err := c.dbPool.Prepare(`
		SELECT
			id,
			customer_number,
			note,
			created_at
		FROM customer_note
	`)
	if err != nil {
		c.logger.Errorf("failed to prepare statement: %v", err)
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		c.logger.Errorf("failed to execute statement: %v", err)
		return nil, err
	}

	defer rows.Close()

	var customerNotes []domain.CustomerNote
	for rows.Next() {
		var customerNote domain.CustomerNote
		err := rows.Scan(
			&customerNote.ID,
			&customerNote.CustomerNumber,
			&customerNote.Note,
			&customerNote.CreatedAt,
		)
		if err != nil {
			c.logger.Errorf("failed to fetch data statement: %v", err)
			return nil, err
		}

		customerNotes = append(customerNotes, customerNote)
	}

	return customerNotes, nil
}

func (c customerNoteRepository) GetByCustomerNumber(customerNumber int, ctx context.Context) ([]domain.CustomerNote, error) {
	stmt, err := c.dbPool.Prepare(`
		SELECT
			id,
			customer_number,
			note,
			created_at
		FROM customer_note
		WHERE
		customer_number = $1
	`)
	if err != nil {
		c.logger.Errorf("failed to prepare statement: %v", err)
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, customerNumber)
	if err != nil {
		c.logger.Errorf("failed to execute statement: %v", err)
		return nil, err
	}

	defer rows.Close()

	var customerNotes []domain.CustomerNote
	for rows.Next() {
		var customerNote domain.CustomerNote
		err := rows.Scan(
			&customerNote.ID,
			&customerNote.CustomerNumber,
			&customerNote.Note,
			&customerNote.CreatedAt,
		)
		if err != nil {
			c.logger.Errorf("failed to fetch data statement: %v", err)
			return nil, err
		}

		customerNotes = append(customerNotes, customerNote)
	}

	return customerNotes, nil
}

func (c customerNoteRepository) GetById(id int, ctx context.Context) (domain.CustomerNote, error) {
	stmt, err := c.dbPool.Prepare(`
		SELECT
			id,
			customer_number,
			note,
			created_at
		FROM customer_note
		WHERE
		id = $1
	`)
	if err != nil {
		c.logger.Errorf("failed to prepare statement: %v", err)
		return domain.CustomerNote{}, err
	}

	var customerNote domain.CustomerNote
	err = stmt.QueryRowContext(ctx, id).Scan(
		&customerNote.ID,
		&customerNote.CustomerNumber,
		&customerNote.Note,
		&customerNote.CreatedAt,
	)
	if err != nil {
		c.logger.Errorf("failed to execute statement: %v", err)
		return domain.CustomerNote{}, err
	}

	return customerNote, nil
}

func (c customerNoteRepository) Insert(customerNote *domain.CustomerNote, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	stmt, err := c.dbPool.PrepareContext(ctx, `
		INSERT INTO customer_note(
			id,
			customer_number,
			note,
			created_at) VALUES (
		$1, $2, $3, $4
	)
	`)
	if err != nil {
		message.Message = "Failed to Insert Customer Note"
		message.StatusCode = 500
		c.logger.Errorf("failed to prepare statement: %v", err)
		return message, nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		customerNote.ID,
		customerNote.CustomerNumber,
		customerNote.Note,
		customerNote.CreatedAt,
	)
	if err != nil {
		message.Message = "Failed to Insert Customer Note"
		message.StatusCode = 500
		c.logger.Errorf("failed to execute statement: %v", err)
		return message, nil
	}
	message.Message = fmt.Sprintf("Succes Insert Customer Note with id %d", customerNote.ID)
	message.StatusCode = 200
	return message, nil
}

func (c customerNoteRepository) Update(customerNote *domain.CustomerNote, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	stmt, err := c.dbPool.PrepareContext(ctx, `
		UPDATE customer_note SET
			customer_number = $2,
			note = $3,
			created_at = $4
		WHERE 
			id = $1
	`)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Update id %d", customerNote.ID)
		message.StatusCode = 500
		c.logger.Errorf("failed to prepare statement: %v", err)
		return message, nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		customerNote.ID,
		customerNote.CustomerNumber,
		customerNote.Note,
		customerNote.CreatedAt,
	)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Update id %d", customerNote.ID)
		message.StatusCode = 500
		c.logger.Errorf("failed to execute statement: %v", err)
		return message, nil
	}
	message.Message = "Succes Update"
	message.StatusCode = 200
	return message, nil
}

func (c customerNoteRepository) DeleteById(id int, ctx context.Context) (domain.Response, error) {
	var message domain.Response
	stmt, err := c.dbPool.Prepare(`
		DELETE
		FROM customer_note
		WHERE id = $1
	`)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Delete id %d", id)
		message.StatusCode = 500
		c.logger.Errorf("failed to prepare statement: %v", err)
		return message, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		message.Message = fmt.Sprintf("Failed Delete id %d", id)
		message.StatusCode = 500
		c.logger.Errorf("failed to execute statement: %v", err)
		return message, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		message.Message = fmt.Sprintf("Could not determine rows affected for id %d", id)
		message.StatusCode = 500
		c.logger.Errorf("failed to get rows affected: %v", err)
		return message, err
	}
	if rowsAffected == 0 {
		message.Message = fmt.Sprintf("No customer note found with id %d", id)
		message.StatusCode = 500
		return message, sql.ErrNoRows
	}
	message.StatusCode = 200
	message.Message = "Succes Delete!!"
	return message, nil
}

func NewCustomerNoteRepository(db *sql.DB, log *logrus.Logger) domain.CustomerNoteRepository {
	return &customerNoteRepository{
		dbPool: db,
		logger: log,
	}
}
