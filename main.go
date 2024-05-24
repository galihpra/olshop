package main

import (
	"olshop/config"
	"olshop/features/users/handler"
	"olshop/features/users/repository"
	"olshop/features/users/service"
	"olshop/helpers/encrypt"
	"olshop/routes"
	"olshop/utilities/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var dbConfig = new(config.DatabaseMysql)
	if err := dbConfig.LoadFromEnv(); err != nil {
		panic(err)
	}

	dbConnection, err := database.MysqlInit(*dbConfig)
	if err != nil {
		panic(err)
	}

	if err := database.MysqlMigrate(dbConnection); err != nil {
		panic(err)
	}

	var jwtConfig = new(config.JWT)
	if err := jwtConfig.LoadFromEnv(); err != nil {
		panic(err)
	}

	enc := encrypt.New()
	userRepository := repository.NewUserRepository(dbConnection)
	userService := service.New(userRepository, enc)
	userHandler := handler.NewUserHandler(userService, *jwtConfig)

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	route := routes.Routes{
		JWTKey:      jwtConfig.Secret,
		Server:      app,
		UserHandler: userHandler,
	}

	route.InitRouter()

	app.Logger.Fatal(app.Start(":8000"))

}
