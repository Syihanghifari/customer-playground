package app

import (
	"context"
	"customer-playground/database"
	"customer-playground/domain"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	delivery_customer "customer-playground/services/customer/delivery"
	repository_customer "customer-playground/services/customer/repository"
	usecase_customer "customer-playground/services/customer/usecase"
	delivery_customernote "customer-playground/services/customernote/delivery"
	repository_customernote "customer-playground/services/customernote/repository"
	usecase_customernote "customer-playground/services/customernote/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run() {
	initConfig()
	logger := initLogger()
	dbPool, err := initDatabase()
	if err != nil {
		logger.Fatalf("%s: %v", "Error on connect to database", err)
	}
	customerNoteUseCase, customerUseCase := initService(dbPool, logger)
	initHandler(customerNoteUseCase, customerUseCase, logger)
}
func initConfig() {
	viper.SetConfigType("toml")

	viper.AddConfigPath(".")
	viper.SetConfigName(".config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	return logger
}

func initDatabase() (*sql.DB, error) {
	dbConn := database.DatabaseConnector{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.name"),
		SSLMode:  viper.GetString("database.sslmode"),
	}

	dbPool, err := dbConn.Connect()
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func initService(dbPool *sql.DB, logger *logrus.Logger) (domain.CustomerNoteUseCase, domain.CustomerUseCase) {
	customerNoteRepository := repository_customernote.NewCustomerNoteRepository(dbPool, logger)
	customerNoteUseCase := usecase_customernote.NewCustomerNoteUseCase(customerNoteRepository, logger)
	customerRepository := repository_customer.NewCustomerRepository(dbPool, logger)
	customerUseCase := usecase_customer.NewCustomerUseCase(customerRepository, logger)
	return customerNoteUseCase, customerUseCase
}

func initHandler(customerNoteUseCase domain.CustomerNoteUseCase, customerUseCase domain.CustomerUseCase, logger *logrus.Logger) {
	ctx := context.Background()

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	http.Handle("/", r)

	// Swagger endpoint
	r.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	delivery_customernote.NewCustomerNoteHandler(r, customerNoteUseCase, logger)
	delivery_customer.NewCustomerHandler(r, customerUseCase, logger)

	srv := &http.Server{
		Addr:         fmt.Sprintf(`:%d`, viper.GetInt("app.port")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
