package main

import (
	"blog/db"
	"blog/middleware"
	"blog/route"
	"context"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(fmt.Errorf("Error when loading .env: %v", err.Error()))
	// }

	dbCon := db.Connect(os.Getenv("MONGODB_CONNECTTION_URI"))
	defer dbCon.Disconnect(context.TODO())

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")

	router.Use(cors.New(corsConfig))
	router.Use(middleware.ErrorHandler)

	route.ConfigRouteAuth(router, dbCon)
	route.ConfigRouteUser(router, dbCon)
	route.ConfigRouteTopic(router, dbCon)
	route.ConfigRoutePost(router, dbCon)

	router.Run(":8080")
}
