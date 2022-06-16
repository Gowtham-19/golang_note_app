package main

import (
	"fmt"

	"github.com/Gowtham-19/note_server/configs"
	"github.com/Gowtham-19/note_server/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Notes Server")
	//initializing router
	router := gin.Default()
	//importing database config
	configs.ConnectDB()
	//Adding user defined routes
	routes.UserRouter(router)
	//starting server
	router.Run("localhost:8000")
}
