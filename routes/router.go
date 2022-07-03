package routes

import (
	"github.com/Gowtham-19/note_golang_server/controller"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {

	router.GET("/getAllNotes", controller.GetAll_Notes)
	router.POST("/createNote", controller.Create_Note)
	router.PUT("/updateNote", controller.Update_Note)
	router.GET("/deleteNote/:id", controller.Delete_Note)
	router.POST("/filterNotes", controller.Filter_Notes)

}
