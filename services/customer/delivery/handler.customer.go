package delivery_customer

import (
	"customer-playground/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CustomerHandler struct {
	customerUseCase domain.CustomerUseCase
	logger          *logrus.Logger
}

func NewCustomerHandler(r *gin.Engine, c domain.CustomerUseCase, l *logrus.Logger) *gin.Engine {
	handler := &CustomerHandler{customerUseCase: c, logger: l}

	r.GET("/customer", handler.HandlerGetAllCustomer)
	r.GET("/customer/:customer_number", handler.HandlerGetCustomerByNumber)
	r.POST("/customer", handler.HandlerInsertCustomer)
	r.PUT("/customer", handler.HandlerUpdateCustomer)
	r.DELETE("/customer/:customer_number", handler.HandlerDeleteCustomerByNumber)

	return r
}

// HandlerGetAllCustomer godoc
// @Summary Get all customers
// @Description Retrieves all customers
// @Tags customers
// @Produce json
// @Success 200 {array} domain.Customer
// @Failure 500 {object} domain.ErrorResponse
// @Router /customer [get]
func (c *CustomerHandler) HandlerGetAllCustomer(ctx *gin.Context) {
	customers, err := c.customerUseCase.GetAll(ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerHandler/HandlerGetAllCustomer", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, customers)
	return
}

// HandlerGetCustomerByNumber godoc
// @Summary Get customer by number
// @Description Retrieves a customer by their customer number
// @Tags customers
// @Produce json
// @Param customer_number path int true "Customer Number"
// @Success 200 {object} domain.Customer
// @Failure 500 {object} domain.ErrorResponse
// @Router /customer/{customer_number} [get]
func (c *CustomerHandler) HandlerGetCustomerByNumber(ctx *gin.Context) {
	customerNumber, err := strconv.Atoi(ctx.Param("customer_number"))
	customer, err := c.customerUseCase.GetByCustomerNumber(customerNumber, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerHandler/HandlerGetCustomerByNumber", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, customer)
	return
}

// HandlerInsertCustomer godoc
// @Summary Insert new customer
// @Description Adds a new customer to the database
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body domain.Customer true "Customer payload"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /customer [post]
func (c *CustomerHandler) HandlerInsertCustomer(ctx *gin.Context) {
	var customer domain.Customer
	err := ctx.Bind(&customer)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerHandler/HandlerInsertCustomer/ParseBodyData", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	message, err := c.customerUseCase.Insert(&customer, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerHandler/HandlerInsertCustomer/Insert", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if message.StatusCode != 200 {
		ctx.JSON(http.StatusBadRequest, message)
		return
	}
	ctx.JSON(http.StatusOK, message)
	return
}

// HandlerUpdateCustomer godoc
// @Summary Update customer
// @Description Updates an existing customer by customer number or ID
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body domain.Customer true "Customer payload"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /customer [put]
func (c *CustomerHandler) HandlerUpdateCustomer(ctx *gin.Context) {
	var customer domain.Customer
	err := ctx.Bind(&customer)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerHandler/HandlerUpdateCustomer/ParseBodyData", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	message, err := c.customerUseCase.Update(&customer, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerHandler/HandlerUpdateCustomer/Update", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if message.StatusCode != 200 {
		ctx.JSON(http.StatusBadRequest, message)
		return
	}
	ctx.JSON(http.StatusOK, message)
	return
}

// HandlerDeleteCustomerByNumber godoc
// @Summary Delete customer by number
// @Description Deletes a customer based on customer number
// @Tags customers
// @Produce json
// @Param customer_number path int true "Customer Number"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /customer/{customer_number} [delete]
func (c *CustomerHandler) HandlerDeleteCustomerByNumber(ctx *gin.Context) {
	customerNumber, err := strconv.Atoi(ctx.Param("customer_number"))
	message, err := c.customerUseCase.DeleteByCustomerNumber(customerNumber, ctx)
	if err != nil {
		c.logger.Errorf("%s : %v", "CustomerHandler/HandlerDeleteCustomerByNumber/Delete", err)
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
