package routes

import (
	"task5-pbi/controllers"
	"task5-pbi/middlewares"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	userRoute := r.Group("/users")
	userRoute.POST("/register", controllers.Register)
	userRoute.POST("/login", controllers.Login)
	userRoute.Use(middlewares.Auth())
	userRoute.PUT("/:userId", controllers.UpdateUserById)
	userRoute.DELETE("/:userId", controllers.DeleteUserById)

	photoRoute := r.Group("/photos")
	photoRoute.GET("/", controllers.GetAllPhotos)
	photoRoute.Use(middlewares.Auth())
	photoRoute.POST("/", controllers.AddPhoto)
	photoRoute.PUT("/:photoId", controllers.UpdatePhotoById)
	photoRoute.DELETE("/:photoId", controllers.DeletePhotoById)

	r.Run()
}
