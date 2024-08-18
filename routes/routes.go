package routes

import (
	"github.com/gonext-tech/internal/handlers"
	"github.com/gonext-tech/internal/manager"
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, store *gorm.DB, projectStores []models.ProjectsDB) {
	userService := services.NewUserServices(models.User{}, store)
	bloodService := services.NewBloodService(models.BloodType{}, store)
	cityService := services.NewCityService(models.City{}, store)
	subscriptionService := services.NewSubscriptionService(projectStores)
	membershipService := services.NewMembershipService(projectStores)
	shopService := services.NewShopService(projectStores)
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
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, projectService)
	membershipHandler := handlers.NewMembershipHandler(membershipService, projectService)
	ticketHandler := handlers.NewTicketHandler(ticketService, projectService, notificationService)
	shopHandler := handlers.NewShopHandler(shopService, projectService)
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

	// membership
	protectedGroup.GET("shop", shopHandler.ListPage)
	protectedGroup.GET("shop/search", shopHandler.SearchUser)
	protectedGroup.GET("shop/view/:id/:name", shopHandler.ViewPage)
	protectedGroup.GET("shop/create", shopHandler.CreatePage)
	protectedGroup.POST("shop/create", shopHandler.CreateHandler)
	protectedGroup.GET("shop/edit/:id/:name", shopHandler.UpdatePage)
	protectedGroup.POST("shop/edit/:id/:name", shopHandler.UpdateHandler)
	protectedGroup.DELETE("shop/:id/:name", shopHandler.DeleteHandler)

	// membership
	protectedGroup.GET("membership", membershipHandler.ListPage)
	protectedGroup.GET("membership/fetch", membershipHandler.Fetch)
	protectedGroup.GET("membership/view/:id/:name", membershipHandler.ViewPage)
	protectedGroup.GET("membership/create", membershipHandler.CreatePage)
	protectedGroup.POST("membership/create", membershipHandler.CreateHandler)
	protectedGroup.GET("membership/edit/:id/:name", membershipHandler.UpdatePage)
	protectedGroup.POST("membership/edit/:id/:name", membershipHandler.UpdateHandler)
	protectedGroup.DELETE("membership/:id/:name", membershipHandler.DeleteHandler)

	// subscrtiption
	protectedGroup.GET("subscription", subscriptionHandler.ListPage)
	//protectedGroup.GET("subscription/fetch", subscriptionHandler.Fetch)
	protectedGroup.GET("subscription/view/:id/:name", subscriptionHandler.ViewPage)
	protectedGroup.GET("subscription/create", subscriptionHandler.CreatePage)
	protectedGroup.POST("subscription/create", subscriptionHandler.CreateHandler)
	protectedGroup.GET("subscription/edit/:id/:name", subscriptionHandler.UpdatePage)
	protectedGroup.POST("subscription/edit/:id/:name", subscriptionHandler.UpdateHandler)
	protectedGroup.DELETE("subscription/:id/:name", subscriptionHandler.DeleteHandler)

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
