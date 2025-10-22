package delivery_customernote

import (
	"customer-playground/domain"
	"net/http"
	"strconv"

	_ "customer-playground/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CustomerNoteHandler struct {
	customerNoteUseCase domain.CustomerNoteUseCase
	logger              *logrus.Logger
}

func NewCustomerNoteHandler(r *gin.Engine, c domain.CustomerNoteUseCase, l *logrus.Logger) *gin.Engine {
	handler := &CustomerNoteHandler{customerNoteUseCase: c, logger: l}

	r.GET("/customer-note/get-all", handler.HandlerGetAllCustomerNote)
	r.GET("/customer-note/get-by-customer-number/:customer_number", handler.HandlerGetByCustomerNumberCustomerNote)
	r.GET("/customer-note/get-by-id/:id", handler.HandlerGetByIdCustomerNote)
	r.POST("/customer-note", handler.HandlerInsertCustomerNote)
	r.PUT("/customer-note", handler.HandlerUpdateCustomerNote)
	r.DELETE("/customer-note/:id", handler.HandlerDeleteCustomerNoteById)

	return r
}

// HandlerGetAllCustomerNote godoc
// @Summary Get all customer notes
// @Description Retrieves all customer notes from the system
// @Tags customer-note
// @Produce json
// @Success 200 {array} domain.CustomerNote "List of customer notes"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /customer-note/get-all [get]
func (c *CustomerNoteHandler) HandlerGetAllCustomerNote(ctx *gin.Context) {
	customerNotes, err := c.customerNoteUseCase.GetAll(ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerGetAllCustomerNote", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, customerNotes)
	return
}

// HandlerGetByCustomerNumberCustomerNote godoc
// @Summary Get customer notes by customer number
// @Description Retrieves all notes associated with a given customer number
// @Tags customer-note
// @Produce json
// @Param customer_number path int true "Customer Number"
// @Success 200 {array} domain.CustomerNote "Customer notes for the customer number"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /customer-note/get-by-customer-number/{customer_number} [get]
func (c *CustomerNoteHandler) HandlerGetByCustomerNumberCustomerNote(ctx *gin.Context) {
	customerNumber, err := strconv.Atoi(ctx.Param("customer_number"))
	customerNotes, err := c.customerNoteUseCase.GetByCustomerNumber(customerNumber, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerGetByCustomerNumberCustomerNote", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, customerNotes)
	return
}

// HandlerGetByIdCustomerNote godoc
// @Summary Get a customer note by ID
// @Description Retrieves a customer note by its unique ID
// @Tags customer-note
// @Produce json
// @Param id path int true "Customer Note ID"
// @Success 200 {object} domain.CustomerNote "Customer note"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /customer-note/get-by-id/{id} [get]
func (c *CustomerNoteHandler) HandlerGetByIdCustomerNote(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	customerNote, err := c.customerNoteUseCase.GetById(id, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerGetByIdCustomerNote", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, customerNote)
	return
}

// HandlerInsertCustomerNote godoc
// @Summary Create a new customer note
// @Description Inserts a new customer note into the system
// @Tags customer-note
// @Accept json
// @Produce json
// @Param customerNote body domain.CustomerNote true "Customer Note Payload"
// @Success 200 {object} domain.Response "Insert result"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /customer-note [post]
func (c *CustomerNoteHandler) HandlerInsertCustomerNote(ctx *gin.Context) {
	var customerNote domain.CustomerNote
	err := ctx.Bind(&customerNote)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerInsertCustomerNote/ParseBodyData", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	message, err := c.customerNoteUseCase.Insert(&customerNote, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerInsertCustomerNote/Insert", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if message.StatusCode != 200 {
		ctx.JSON(http.StatusBadRequest, message)
		return
	}
	ctx.JSON(http.StatusOK, message)
	return
}

// HandlerUpdateCustomerNote godoc
// @Summary Update a customer note
// @Description Updates an existing customer note
// @Tags customer-note
// @Accept json
// @Produce json
// @Param customerNote body domain.CustomerNote true "Customer Note Payload"
// @Success 200 {object} domain.Response "Update result"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /customer-note [put]
func (c *CustomerNoteHandler) HandlerUpdateCustomerNote(ctx *gin.Context) {
	var customerNote domain.CustomerNote
	err := ctx.Bind(&customerNote)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerUpdateCustomerNote/ParseBodyData", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	message, err := c.customerNoteUseCase.Update(&customerNote, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerUpdateCustomerNote/Update", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if message.StatusCode != 200 {
		ctx.JSON(http.StatusBadRequest, message)
		return
	}
	ctx.JSON(http.StatusOK, message)
	return
}

// HandlerDeleteCustomerNoteById godoc
// @Summary Delete a customer note by ID
// @Description Deletes a customer note by its ID
// @Tags customer-note
// @Produce json
// @Param id path int true "Customer Note ID"
// @Success 200 {object} domain.Response "Delete result"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /customer-note/{id} [delete]
func (c *CustomerNoteHandler) HandlerDeleteCustomerNoteById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	message, err := c.customerNoteUseCase.DeleteById(id, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerNoteHandler/HandlerDeleteCustomerNoteById/Delete", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if message.StatusCode != 200 {
		ctx.JSON(http.StatusBadRequest, message)
		return
	}
	ctx.JSON(http.StatusOK, message)
	return
}
