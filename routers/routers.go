package router

import (
	"jobApps/authentication"
	"jobApps/handlers"
	"jobApps/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func Router(conn *pgx.Conn) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	query := database.New(conn)

	handler := handlers.ControllerInstance(query)

	//signup
	router.POST("/signup", handler.SignUp)

	//login
	router.POST("/login", handler.Login)

	// Career
	router.POST("/createcareer", authentication.AuthMiddleware(), handler.CreateCareer)
	router.GET("/getcareerdetail/:id", authentication.AuthMiddleware(), handler.GetCareerByJobId)
	router.GET("/get-all-career-details", authentication.AuthMiddleware(), handler.GetAllCareers)
	router.PUT("/updatecareer/:id", authentication.AuthMiddleware(), handler.UpdateCareerById)
	router.DELETE("/deletecareer/:id", authentication.AuthMiddleware(), handler.DeleteCareerById)

	// Profile
	router.POST("/createprofile", authentication.AuthMiddleware(), handler.CreateProfile)
	router.GET("/getprofile/:id", authentication.AuthMiddleware(), handler.GetProfileById)
	router.GET("/get-all-profile-details", authentication.AuthMiddleware(), handler.GetAllProfiles)
	router.DELETE("/delete-profile/:id", authentication.AuthMiddleware(), handler.DeleteProfileById)
	router.PUT("/update-profile/:id", authentication.AuthMiddleware(), handler.UpdateProfileById)

	//User
	router.GET("/get-all-users-email", authentication.AuthMiddleware(), handler.GetAllUsersEmail)

	//router
	router.Run("localhost:8080")
}
