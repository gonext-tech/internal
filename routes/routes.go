package routes

import (
	"github.com/gonext-tech/internal/handlers"
	//"github.com/gonext-tech/internal/manager"
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, store *gorm.DB, projectStores []models.ProjectsDB) {
	userService := services.NewUserServices(models.User{}, store)
	subscriptionService := services.NewSubscriptionService(projectStores)
	membershipService := services.NewMembershipService(projectStores)
	shopService := services.NewShopService(projectStores)
	projectService := services.NewProjectService(models.Project{}, store)
	authHandler := handlers.NewAuthHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, projectService)
	membershipHandler := handlers.NewMembershipHandler(membershipService, projectService)
	shopHandler := handlers.NewShopHandler(shopService, projectService)

	// --> UNPRTECTED ROUTES <--
	e.GET("/", authHandler.HomeHandler)
	e.GET("/login", authHandler.LoginHandler)
	e.POST("/login", authHandler.LoginHandler)
	e.GET("/register", authHandler.RegisterHandler)
	e.POST("/register", authHandler.RegisterHandler)

	protectedGroup := e.Group("/", authHandler.AuthMiddleware)

	// --> WEBSOCKET ROUTES <--
	//protectedGroup.GET("ws/ticket", manager.Connect)

	// --> PROJECT ROUTES <--
	protectedGroup.GET("project", projectHandler.ListPage)
	protectedGroup.GET("project/view", projectHandler.ViewPage)
	protectedGroup.GET("project/create", projectHandler.CreatePage)
	protectedGroup.POST("project/create", projectHandler.CreateHandler)
	protectedGroup.GET("project/edit/:id", projectHandler.UpdatePage)
	protectedGroup.POST("project/edit/:id", projectHandler.UpdateHandler)
	protectedGroup.DELETE("project/:id", projectHandler.DeleteHandler)

	// --> SHOP ROUTES <--
	protectedGroup.GET("shop", shopHandler.ListPage)
	protectedGroup.GET("shop/search", shopHandler.SearchUser)
	protectedGroup.GET("shop/view/:id/:name", shopHandler.ViewPage)
	protectedGroup.GET("shop/create", shopHandler.CreatePage)
	protectedGroup.POST("shop/create", shopHandler.CreateHandler)
	protectedGroup.GET("shop/edit/:id/:name", shopHandler.UpdatePage)
	protectedGroup.POST("shop/edit/:id/:name", shopHandler.UpdateHandler)
	protectedGroup.DELETE("shop/:id/:name", shopHandler.DeleteHandler)

	// --> MEMBERSHIP ROUTES <--
	protectedGroup.GET("membership", membershipHandler.ListPage)
	protectedGroup.GET("membership/fetch", membershipHandler.Fetch)
	protectedGroup.GET("membership/view/:id/:name", membershipHandler.ViewPage)
	protectedGroup.GET("membership/create", membershipHandler.CreatePage)
	protectedGroup.POST("membership/create", membershipHandler.CreateHandler)
	protectedGroup.GET("membership/edit/:id/:name", membershipHandler.UpdatePage)
	protectedGroup.POST("membership/edit/:id/:name", membershipHandler.UpdateHandler)
	protectedGroup.DELETE("membership/:id/:name", membershipHandler.DeleteHandler)

	// --> SUBSCRIPTION ROUTES <--
	protectedGroup.GET("subscription", subscriptionHandler.ListPage)
	protectedGroup.GET("subscription/view/:id/:name", subscriptionHandler.ViewPage)
	protectedGroup.GET("subscription/create", subscriptionHandler.CreatePage)
	protectedGroup.POST("subscription/create", subscriptionHandler.CreateHandler)
	protectedGroup.GET("subscription/edit/:id/:name", subscriptionHandler.UpdatePage)
	protectedGroup.POST("subscription/edit/:id/:name", subscriptionHandler.UpdateHandler)
	protectedGroup.DELETE("subscription/:id/:name", subscriptionHandler.DeleteHandler)

	protectedGroup.GET("user/search", authHandler.SearchUser)

}
