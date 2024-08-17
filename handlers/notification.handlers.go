package handlers

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/components"
	"github.com/gonext-tech/internal/views/notification_views"
)

type NotificationService interface {
	GetALL(limit, page int, orderBy, sortBy string) ([]models.Notification, models.Meta, error)
	GetNew() ([]models.Notification, error)
	Read() error
	Create(ticket models.Ticket, label string) error
}

type NotificationHandler struct {
	NotificationServices NotificationService
}

func NewNotificationHandler(ns NotificationService) *NotificationHandler {
	return &NotificationHandler{
		NotificationServices: ns,
	}
}
func (nh *NotificationHandler) SearchPage(c echo.Context) error {
	notifications, err := nh.NotificationServices.GetNew()
	if err != nil {
		isError = false
		errorMsg = "can't fetch projects"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	return renderView(c, components.NotificationItems(notifications))
}
func (nh *NotificationHandler) ListPage(c echo.Context) error {
	isError = false
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}
	orderBy := c.QueryParam("orderBy")
	if orderBy == "" {
		orderBy = "desc"
	}
	sortBy := c.QueryParam("sortBy")
	if sortBy == "" {
		sortBy = "id"
	}
	notifications, meta, err := nh.NotificationServices.GetALL(limit, page, orderBy, sortBy)
	if err != nil {
		isError = false
		errorMsg = "can't fetch notifications"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}
	titlePage := fmt.Sprintf(
		"Notifications (%d)", meta.TotalCount)
	return renderView(c, notification_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		notification_views.List(titlePage, notifications, meta),
	))
}

func (nh *NotificationHandler) Create(c echo.Context) error {
	return nil
}
