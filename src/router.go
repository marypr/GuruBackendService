package src

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const (
	apiRoute = "/api/v1"
	port     = "8080"
)

//Router is a router variable
var Router *gin.Engine

//Shutdown is a channel to shutdown the router in runtime
var Shutdown chan int

//Start is a function that starts server and initializes routes
func Start() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()
	Router.Use(gin.Logger())
	Router.Use(gin.Recovery())

	initUserProfileRoutes()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: Router,
	}

	go func() {
		// Starting server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	Shutdown = make(chan int)
	<-Shutdown
	log.Println("Shutdown Server ...")
}

func initUserProfileRoutes() {
	userService := NewUserService(NewPostgresUsersRepo(Connection))

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-Shutdown:
				return
			case <-ticker.C:
				userService.updateUsers()
			}
		}
	}()

	userRoutes := Router.Group(apiRoute)
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
