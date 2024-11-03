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

	// --> SERVICES INIT <--
	userService := services.NewUserServices(models.User{}, store)
	referalService := services.NewReferalService(models.Referal{}, store)
	subscriptionService := services.NewSubscriptionService(projectStores)
	membershipService := services.NewMembershipService(projectStores)
	shopService := services.NewShopService(projectStores)
	projectService := services.NewProjectService(models.Project{}, store)
	statsService := services.NewStatisticcServices(models.Stats{}, store)
	invoiceService := services.NewInvoiceService(models.Invoice{}, store)
	customerService := services.NewCustomerService(projectStores)
	uploadService := services.NewUploadServices(store)
	appointmentService := services.NewAppointmentService(projectStores)

	// --> HANDLERS INIT <--
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	referalHandler := handlers.NewReferalHandler(referalService, uploadService)
	authHandler := handlers.NewAuthHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService, uploadService)
	statsHandler := handlers.NewStatsHandler(statsService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService, projectService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, projectService, membershipService, shopService, statsService)
	membershipHandler := handlers.NewMembershipHandler(membershipService, projectService)
	shopHandler := handlers.NewShopHandler(shopService, projectService, membershipService, uploadService)
	automationHandler := handlers.NewAutomationHandler(projectStores)
	customerHandler := handlers.NewCustomerHandler(customerService, projectService, uploadService)
	_ = handlers.NewUploadHandler(uploadService, projectService)

	// --> UNPRTECTED ROUTES <--
	e.GET("/", authHandler.HomeHandler)
	e.GET("/login", authHandler.LoginHandler)
	e.POST("/login", authHandler.LoginHandler)
	e.GET("/register", authHandler.RegisterHandler)
	e.POST("/register", authHandler.RegisterHandler)
	e.POST("/logout", authHandler.LogoutHandler)

	// --> AUTOMATION <--
	e.GET("/api/send/wp", automationHandler.GetAppointments)
	e.PUT("/api/send/wp/:id", automationHandler.UpdateAppointment)

	protectedGroup := e.Group("/", authHandler.AuthMiddleware)

	// --> WEBSOCKET ROUTES <--
	//protectedGroup.GET("ws/ticket", manager.Connect)

	protectedGroup.GET("dashboard", statsHandler.DashboardPage)

	// --> PROJECT ROUTES <--
	protectedGroup.GET("project", projectHandler.ListPage)
	protectedGroup.GET("project/view", projectHandler.ViewPage)
	protectedGroup.GET("project/create", projectHandler.CreatePage)
	protectedGroup.POST("project/create", projectHandler.CreateHandler)
	protectedGroup.GET("project/edit/:id", projectHandler.UpdatePage)
	protectedGroup.POST("project/edit/:id", projectHandler.UpdateHandler)
	protectedGroup.DELETE("project/:id", projectHandler.DeleteHandler)

	// --> PROJECT ROUTES <--
	protectedGroup.GET("invoice", invoiceHandler.ListPage)
	protectedGroup.GET("invoice/view", invoiceHandler.ViewPage)
	protectedGroup.GET("invoice/create", invoiceHandler.CreatePage)
	protectedGroup.POST("invoice/create", invoiceHandler.CreateHandler)
	protectedGroup.GET("invoice/edit/:id", invoiceHandler.UpdatePage)
	protectedGroup.POST("invoice/edit/:id", invoiceHandler.UpdateHandler)
	protectedGroup.DELETE("invoice/:id", invoiceHandler.DeleteHandler)

	// --> Referal ROUTES <--
	protectedGroup.GET("referal", referalHandler.ListPage)
	protectedGroup.GET("referal/view", referalHandler.ViewPage)
	protectedGroup.GET("referal/create", referalHandler.CreatePage)
	protectedGroup.POST("referal/create", referalHandler.CreateHandler)
	protectedGroup.GET("referal/edit/:id", referalHandler.UpdatePage)
	protectedGroup.POST("referal/edit/:id", referalHandler.UpdateHandler)
	protectedGroup.DELETE("referal/:id", referalHandler.DeleteHandler)

	// --> SHOP ROUTES <--
	protectedGroup.GET("shop", shopHandler.ListPage)
	protectedGroup.GET("shop/appointment/:id", appointmentHandler.ListPage)
	protectedGroup.GET("shop/search", shopHandler.SearchUser)
	protectedGroup.GET("shop/view/:id/:name", shopHandler.ViewPage)
	protectedGroup.GET("shop/create", shopHandler.CreatePage)
	protectedGroup.POST("shop/create", shopHandler.CreateHandler)
	protectedGroup.GET("shop/edit/:id/:name", shopHandler.UpdatePage)
	protectedGroup.GET("shop/subscription/:id/:name", subscriptionHandler.ShopListPage)
	protectedGroup.POST("shop/edit/:id/:name", shopHandler.UpdateHandler)
	protectedGroup.DELETE("shop/:id/:name", shopHandler.DeleteHandler)

	// --> CUSTOMER ROUTES <--
	protectedGroup.GET("customer", customerHandler.ListPage)
	protectedGroup.GET("customer/search", customerHandler.SearchUser)
	protectedGroup.GET("customer/view/:id/:name", customerHandler.ViewPage)
	protectedGroup.GET("customer/create", customerHandler.CreatePage)
	protectedGroup.POST("customer/create", customerHandler.CreateHandler)
	protectedGroup.GET("customer/edit/:id/:name", customerHandler.UpdatePage)
	protectedGroup.POST("customer/edit/:id/:name", customerHandler.UpdateHandler)
	protectedGroup.DELETE("customer/:id/:name", customerHandler.DeleteHandler)

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

	// --> SUBSCRIPTION ROUTES <--

	protectedGroup.GET("user/search", authHandler.SearchUser)

}
