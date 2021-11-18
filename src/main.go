package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"golang_bank_demo/src/api"
	"golang_bank_demo/src/config"
	"golang_bank_demo/src/postgres"
	"golang_bank_demo/src/service"
	"golang_bank_demo/src/storage"
	"log"
	"net/http"
)

func main() {
	var appConfig config.AppConfig
	if err := cleanenv.ReadConfig("config.yaml", &appConfig); err != nil {
		log.Fatal(err)
	} else if pgClient, err := postgres.CreateClient(&appConfig); err != nil {
		log.Fatal(err)
	} else if err = postgres.SetUp(pgClient); err != nil {
		log.Fatal(err)
	} else {
		accountStorage := storage.NewPostgresAccountStorage(pgClient)
		accountService := service.NewAccountService(accountStorage)
		authService := service.NewStubAuthenticationService()
		auth := api.NewAuthenticatedApi(authService)
		accountApi := api.NewAccountApi(accountService, auth)

		done := make(chan bool)
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Port), accountApi.Router()))
		}()
		log.Printf("Server started on port %v", appConfig.Port)
		<-done
	}
}
