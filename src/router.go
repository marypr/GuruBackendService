package src

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	apiRoute = "/api/v1"
	port     = ":8080"
	tenSec   = 10
)

//router is a variable for gin router
var router *gin.Engine

//shutdown is a channel to shutdown the router in runtime
var shutdown chan int

//Start is a function that starts server and initializes routes
func Start() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	initUserProfileRoutes()

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		// Starting server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	shutdown = make(chan int)
	<-shutdown
	log.Println("shutdown Server ...")
}

func initUserProfileRoutes() {
	userService := NewUserService(NewPostgresUsersRepo(connection))
	ticker := time.NewTicker(tenSec * time.Second)
	userRoutes := router.Group(apiRoute)

	//saves modified users every 10seconds
	go func() {
		for {
			select {
			case <-shutdown:
				return
			case <-ticker.C:
				userService.updateUsers()
			}
		}
	}()

	{
		// Handle POST requests at /api/v1/user/create
		userRoutes.POST("/user/create", userService.AddUser)
		// Handle POST requests at /api/v1/user/get
		userRoutes.POST("/user/get", userService.GetUser)
		// Handle POST requests at /api/v1/user/deposit
		userRoutes.POST("/user/deposit", userService.AddDeposit)
		// Handle POST requests at /api/v1/transaction
		userRoutes.POST("/transaction", userService.MakeTransaction)
	}
}
