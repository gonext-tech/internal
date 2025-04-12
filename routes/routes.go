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
	referalService := services.NewReferalService(models.Referal{}, store)
	clientService := services.NewClientService(models.Client{}, store)
	membershipService := services.NewMembershipService(projectStores)
	shopService := services.NewShopService(projectStores)
	projectService := services.NewProjectService(models.Project{}, store)
	statsService := services.NewStatisticcServices(models.Stats{}, store)
	invoiceService := services.NewInvoiceService(models.Invoice{}, store)
	uploadService := services.NewUploadServices(store)
	appointmentService := services.NewAppointmentService(projectStores)
	serverService := services.NewServerService(models.MonitoredServer{}, store)
	domainService := services.NewDomainService(models.Domain{}, store)
	adminService := services.NewAdminService(models.Admin{}, store)

	// --> HANDLERS INIT <--
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	referalHandler := handlers.NewReferalHandler(referalService, uploadService)
	authHandler := handlers.NewAuthHandler(adminService)
	projectHandler := handlers.NewProjectHandler(projectService, uploadService, serverService, adminService)
	serverHandler := handlers.NewServerHandler(serverService, uploadService)
	domainHandler := handlers.NewDomainHandler(domainService, serverService)
	statsHandler := handlers.NewStatsHandler(statsService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService, projectService)
	//subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, projectService, membershipService, shopService, statsService)
	membershipHandler := handlers.NewMembershipHandler(membershipService, projectService)
	shopHandler := handlers.NewShopHandler(shopService, projectService, membershipService, uploadService)
	automationHandler := handlers.NewAutomationHandler(projectStores)
	//customerHandler := handlers.NewCustomerHandler(customerService, projectService, uploadService)
	adminHandler := handlers.NewAdminHandler(adminService, uploadService)
	clientHandler := handlers.NewClientHandler(clientService, uploadService)
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

	// --> ADMIN ROUTES <--
	protectedGroup.GET("admin", adminHandler.ListPage)
	protectedGroup.GET("admin/view", adminHandler.ViewPage)
	protectedGroup.GET("admin/create", adminHandler.CreatePage)
	protectedGroup.POST("admin/create", adminHandler.CreateHandler)
	protectedGroup.GET("admin/edit/:id", adminHandler.UpdatePage)
	protectedGroup.POST("admin/edit/:id", adminHandler.UpdateHandler)
	protectedGroup.DELETE("admin/:id", adminHandler.DeleteHandler)

	// --> INVOICE ROUTES <--
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
	protectedGroup.POST("shop/edit/:id/:name", shopHandler.UpdateHandler)
	protectedGroup.DELETE("shop/:id/:name", shopHandler.DeleteHandler)

	// --> SERVER ROUTES <--
	protectedGroup.GET("server", serverHandler.ListPage)
	protectedGroup.GET("server/view/:id", serverHandler.ViewPage)
	protectedGroup.GET("server/create", serverHandler.CreatePage)
	protectedGroup.POST("server/create", serverHandler.CreateHandler)
	protectedGroup.GET("server/edit/:id", serverHandler.UpdatePage)
	protectedGroup.POST("server/edit/:id", serverHandler.UpdateHandler)
	protectedGroup.DELETE("server/:id", serverHandler.DeleteHandler)

	// --> CLIENT ROUTES <--
	protectedGroup.GET("client", clientHandler.ListPage)
	protectedGroup.GET("client/view/:id", clientHandler.ViewPage)
	protectedGroup.GET("client/create", clientHandler.CreatePage)
	protectedGroup.POST("client/create", clientHandler.CreateHandler)
	protectedGroup.GET("client/edit/:id", clientHandler.UpdatePage)
	protectedGroup.POST("client/edit/:id", clientHandler.UpdateHandler)
	protectedGroup.DELETE("client/:id", clientHandler.DeleteHandler)

	// --> DOMAIN ROUTES <--
	protectedGroup.GET("domain", domainHandler.ListPage)
	protectedGroup.GET("domain/view/:id", domainHandler.ViewPage)
	protectedGroup.GET("domain/create", domainHandler.CreatePage)
	protectedGroup.POST("domain/create", domainHandler.CreateHandler)
	protectedGroup.GET("domain/edit/:id", domainHandler.UpdatePage)
	protectedGroup.POST("domain/edit/:id", domainHandler.UpdateHandler)
	protectedGroup.DELETE("domain/:id", domainHandler.DeleteHandler)

	// --> MEMBERSHIP ROUTES <--
	protectedGroup.GET("membership", membershipHandler.ListPage)
	protectedGroup.GET("membership/fetch", membershipHandler.Fetch)
	protectedGroup.GET("membership/view/:id/:name", membershipHandler.ViewPage)
	protectedGroup.GET("membership/create", membershipHandler.CreatePage)
	protectedGroup.POST("membership/create", membershipHandler.CreateHandler)
	protectedGroup.GET("membership/edit/:id/:name", membershipHandler.UpdatePage)
	protectedGroup.POST("membership/edit/:id/:name", membershipHandler.UpdateHandler)
	protectedGroup.DELETE("membership/:id/:name", membershipHandler.DeleteHandler)

	protectedGroup.GET("user/search", authHandler.SearchUser)

}
