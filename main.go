package main

import (
	"olshop/config"
	ph "olshop/features/products/handler"
	pr "olshop/features/products/repository"
	ps "olshop/features/products/service"
	uh "olshop/features/users/handler"
	ur "olshop/features/users/repository"
	us "olshop/features/users/service"
	"olshop/helpers/encrypt"
	"olshop/routes"
	"olshop/utilities/cloudinary"
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

	var cloudinaryConfig = new(config.Cloudinary)
	if err := cloudinaryConfig.LoadFromEnv(); err != nil {
		panic(err)
	}

	cloudinary, err := cloudinary.NewCloudinary(*cloudinaryConfig)
	if err != nil {
		panic(err)
	}

	enc := encrypt.New()
	userRepository := ur.NewUserRepository(dbConnection, cloudinary)
	userService := us.New(userRepository, enc)
	userHandler := uh.NewUserHandler(userService, *jwtConfig)

	productRepository := pr.NewProductRepository(dbConnection, cloudinary)
	productService := ps.NewProductService(productRepository)
	productHandler := ph.NewProductHandler(productService)

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	route := routes.Routes{
		JWTKey:         jwtConfig.Secret,
		Server:         app,
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}

	route.InitRouter()

	app.Logger.Fatal(app.Start(":8000"))

}
