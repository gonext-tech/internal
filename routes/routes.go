package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ramyjaber1/internal/handlers"
	"github.com/ramyjaber1/internal/manager"
	"github.com/ramyjaber1/internal/models"
	"github.com/ramyjaber1/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, store *gorm.DB) {
	userService := services.NewUserServices(models.User{}, store)
	bloodService := services.NewBloodService(models.BloodType{}, store)
	cityService := services.NewCityService(models.City{}, store)
	donorService := services.NewDonorService(models.Donor{}, store)
	projectService := services.NewProjectService(models.Project{}, store)
	ticketService := services.NewTicketService(models.Ticket{}, store)
	notificationService := services.NewNotificationService(models.Notification{}, store)
	authHandler := handlers.NewAuthHandler(userService)
	bloodHandler := handlers.NewBloodHandler(bloodService)
	cityHandler := handlers.NewCityHandler(cityService)
	donorHandler := handlers.NewDonorHandler(donorService, bloodService, cityService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	projectHandler := handlers.NewProjectHandler(projectService)
	ticketHandler := handlers.NewTicketHandler(ticketService, projectService, notificationService)
	e.GET("/", authHandler.HomeHandler)
	e.GET("/login", authHandler.LoginHandler)
	e.POST("/login", authHandler.LoginHandler)
	e.GET("/register", authHandler.RegisterHandler)
	e.POST("/register", authHandler.RegisterHandler)

	e.POST("/api/v1/ticket", ticketHandler.ApiCreate)

	protectedGroup := e.Group("/", authHandler.AuthMiddleware)

	protectedGroup.GET("ws/ticket", manager.Connect)
	// PROJECT
	protectedGroup.GET("project", projectHandler.ListPage)
	protectedGroup.GET("project/view", projectHandler.ViewPage)
	protectedGroup.GET("project/create", projectHandler.CreatePage)
	protectedGroup.POST("project/create", projectHandler.CreateHandler)
	protectedGroup.GET("project/edit/:id", projectHandler.UpdatePage)
	protectedGroup.POST("project/edit/:id", projectHandler.UpdateHandler)
	protectedGroup.DELETE("project/:id", projectHandler.DeleteHandler)

	// Blood-type
	protectedGroup.GET("blood-type", bloodHandler.ListPage)
	protectedGroup.GET("blood-type/view", bloodHandler.ViewPage)
	protectedGroup.GET("blood-type/create", bloodHandler.CreatePage)
	protectedGroup.POST("blood-type/create", bloodHandler.CreateHandler)
	protectedGroup.GET("blood-type/edit/:id", bloodHandler.UpdatePage)
	protectedGroup.POST("blood-type/edit/:id", bloodHandler.UpdateHandler)
	protectedGroup.DELETE("blood-type/:id", bloodHandler.DeleteHandler)

	// city
	protectedGroup.GET("city", cityHandler.ListPage)
	protectedGroup.GET("city/view", cityHandler.ViewPage)
	protectedGroup.GET("city/create", cityHandler.CreatePage)
	protectedGroup.POST("city/create", cityHandler.CreateHandler)
	protectedGroup.GET("city/edit/:id", cityHandler.UpdatePage)
	protectedGroup.POST("city/edit/:id", cityHandler.UpdateHandler)
	protectedGroup.DELETE("city/:id", cityHandler.DeleteHandler)

	// donor
	protectedGroup.GET("donor", donorHandler.ListPage)
	protectedGroup.GET("donor/view", donorHandler.ViewPage)
	protectedGroup.GET("donor/create", donorHandler.CreatePage)
	protectedGroup.POST("donor/create", donorHandler.CreateHandler)
	protectedGroup.GET("donor/edit/:id", donorHandler.UpdatePage)
	protectedGroup.POST("donor/edit/:id", donorHandler.UpdateHandler)
	protectedGroup.DELETE("donor/:id", donorHandler.DeleteHandler)

	protectedGroup.GET("user/search", authHandler.SearchUser)
	protectedGroup.GET("ticket", ticketHandler.ListPage)
	protectedGroup.GET("ticket/search", ticketHandler.Search)
	protectedGroup.GET("ticket/view", ticketHandler.ViewPage)
	protectedGroup.GET("ticket/create", ticketHandler.CreatePage)
	protectedGroup.POST("ticket/create", ticketHandler.CreateHandler)
	protectedGroup.GET("ticket/edit/:id", ticketHandler.UpdatePage)
	protectedGroup.POST("ticket/edit/:id", ticketHandler.UpdateHandler)
	protectedGroup.DELETE("ticket/:id", ticketHandler.DeleteHandler)
	protectedGroup.GET("notification/navbar", notificationHandler.SearchPage)
}
