package main

import (
	"log"
	"net/http"
	"pomodoro-rpg-api/cmd/api/router"
	"pomodoro-rpg-api/infra/db"
	"pomodoro-rpg-api/infra/persistence"
	"pomodoro-rpg-api/infra/service"
	"pomodoro-rpg-api/pkg/config"
	"pomodoro-rpg-api/pkg/logger"
	"pomodoro-rpg-api/presentation/handler"
	"pomodoro-rpg-api/presentation/middleware"
	"pomodoro-rpg-api/usecase"
)

func main() {
	conf := config.NewConfig()
	db, err := db.NewDB(conf.DB)
	if err != nil {
		log.Fatalf("database initialize failed: %v", err)
	}

	logger.Init()

	accRepo := persistence.NewaccountPersistence(db)
	accUsecase := usecase.NewAccountUsecase(accRepo)
	accHandler := handler.NewAccountHandler(accUsecase)

	tr := persistence.NewTimePersistence(db)
	tu := usecase.NewTimeUsecase(accRepo, tr)
	th := handler.NewTimeHandler(tu)

	cognitoService, err := service.NewCognitoService(conf.AWS.ClientID, conf.AWS.ClientSecret, conf.AWS.UserPoolID)
	if err != nil {
		log.Fatalf("cognito initialize failed: %v", err)
	}

	authUsecase := usecase.NewAuthUsecase(cognitoService, accRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	// TODO: middlewareがcognitoServiceに依存するのはイマイチなのでリファクタする
	authenticator := middleware.NewAuthenticator(conf.AWS.UserPoolID, conf.AWS.ClientID, cognitoService)

	deps := router.HandlerDependencies{
		AuthHandler:    authHandler,
		AccountHandler: accHandler,
		TimeHandler:    th,
	}

	r := router.New(deps, authenticator)

	log.Println("🚀 Server is running!")
	http.ListenAndServe(":8080", r)
}
