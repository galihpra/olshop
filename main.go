package main

import (
	"olshop/config"
	ah "olshop/features/addresses/handler"
	ar "olshop/features/addresses/repository"
	as "olshop/features/addresses/service"
	ch "olshop/features/carts/handler"
	cr "olshop/features/carts/repository"
	cs "olshop/features/carts/service"
	ph "olshop/features/products/handler"
	pr "olshop/features/products/repository"
	ps "olshop/features/products/service"
	rh "olshop/features/reviews/handler"
	rr "olshop/features/reviews/repository"
	rs "olshop/features/reviews/service"
	th "olshop/features/transactions/handler"
	tr "olshop/features/transactions/repository"
	ts "olshop/features/transactions/service"
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

	addressRepository := ar.NewAddressRepository(dbConnection)
	addressService := as.NewAddressService(addressRepository)
	addressHandler := ah.NewAddressHandler(addressService, *jwtConfig)

	productRepository := pr.NewProductRepository(dbConnection, cloudinary)
	productService := ps.NewProductService(productRepository)
	productHandler := ph.NewProductHandler(productService)

	reviewRepository := rr.NewReviewRepository(dbConnection)
	reviewService := rs.NewReviewService(reviewRepository)
	reviewHandler := rh.NewReviewHandler(reviewService, *jwtConfig)

	cartRepository := cr.NewCartRepository(dbConnection)
	cartService := cs.NewCartService(cartRepository)
	cartHandler := ch.NewCartHandler(cartService, *jwtConfig)

	transactionRepository := tr.NewTransactionRepository(dbConnection)
	transactionService := ts.NewTransactionService(transactionRepository)
	transactionHandler := th.NewTransactionHandler(transactionService, *jwtConfig)

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	route := routes.Routes{
		JWTKey:             jwtConfig.Secret,
		Server:             app,
		UserHandler:        userHandler,
		AddressHandler:     addressHandler,
		ProductHandler:     productHandler,
		ReviewHandler:      reviewHandler,
		CartHandler:        cartHandler,
		TransactionHandler: transactionHandler,
	}

	route.InitRouter()

	app.Logger.Fatal(app.Start(":8000"))

}
