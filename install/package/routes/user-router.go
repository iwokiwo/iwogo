package routes

import (
	"iwogo/auth"
	"iwogo/middleware"
	"iwogo/modules/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRouter(db *gorm.DB, api *gin.RouterGroup, auth auth.Service) *gin.RouterGroup {

	// injector
	//authService := auth.NewService()
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userController := user.NewUserController(userService, auth)

	// register and login
	api.POST("/register", userController.RegisterUser)
	api.POST("/sessions", userController.Login)
	api.POST("/user/email/checker", userController.ChekEmailAvailability)

	// AFTER LOGIN
	api.POST("/user/change/name", middleware.AuthMiddleware(auth, userService), userController.ChangeNameHandler)
	api.GET("/user/detail", middleware.AuthMiddleware(auth, userService), userController.FetchUser)
	api.GET("/users", middleware.AuthMiddleware(auth, userService), userController.GetAllUsers)
	api.POST("/user/change/password", middleware.AuthMiddleware(auth, userService), userController.ChangePassword)
	api.POST("/user/delete", middleware.AuthMiddleware(auth, userService), userController.DeleteUser)
	api.POST("/user/detail", middleware.AuthMiddleware(auth, userService), userController.GetUserByID)
	api.POST("/user/change", middleware.AuthMiddleware(auth, userService), userController.ChangeDetailHandler)

	return api
}
