package main

import (
	"fmt"
	"os"

	"github.com/Gowtham-19/note_golang_server/configs"
	"github.com/Gowtham-19/note_golang_server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Notes Server")
	//initializing router
	router := gin.Default()
	//setting cors
	router.Use(cors.Default())
	//importing database config
	configs.ConnectDB()
	//Adding user defined routes
	routes.UserRouter(router)
	port := os.Getenv("PORT")
	//starting server
	sever_host := "localhost:" + port
	router.Run(sever_host)
}
