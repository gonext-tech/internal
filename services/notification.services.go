package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/gonext-tech/internal/manager"
	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type NotificationServices struct {
	Notification models.Notification
	DB           *gorm.DB
}

func NewNotificationService(n models.Notification, db *gorm.DB) *NotificationServices {
	return &NotificationServices{
		Notification: n,
		DB:           db,
	}
}

func (ns *NotificationServices) GetALL(limit, page int, orderBy, sortBy string) ([]models.Notification, models.Meta, error) {
	var notifications []models.Notification
	query := ns.DB
	totalQuery := query
	offset := (page - 1) * limit
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&notifications)
	totalRecords := int64(0)
	totalQuery.Model(&ns.Notification).Count(&totalRecords)
	lastPage := int64(0)
	if limit > 0 {
		lastPage = (totalRecords + int64(limit) - 1) / int64(limit)
	}
	meta := models.Meta{
		CurrentPage: page,
		TotalCount:  int(totalRecords),
		LastPage:    int(lastPage),
		Limit:       limit,
	}

	return notifications, meta, nil
}

func (ns *NotificationServices) GetNew() ([]models.Notification, error) {
	var notificationsUnread []models.Notification
	var notificationsRead []models.Notification
	if result := ns.DB.Order("created_at DESC").Where("new = ?", true).Find(&notificationsUnread); result.Error != nil {
		return []models.Notification{}, result.Error
	}
	if result := ns.DB.Order("created_at DESC").Where("new = ?", false).Limit(5).Find(&notificationsRead); result.Error != nil {
		return []models.Notification{}, result.Error
	}
	var combinedNotifications []models.Notification
	combinedNotifications = append(combinedNotifications, notificationsUnread...)
	combinedNotifications = append(combinedNotifications, notificationsRead...)

	return combinedNotifications, nil
}

func (ns *NotificationServices) Read() error {
	if result := ns.DB.Model(&models.Notification{}).Where("new = ?", false).Updates(map[string]interface{}{"read": true}); result.Error != nil {
		return result.Error
	}
	return nil
}

func (ns *NotificationServices) Create(ticket models.Ticket, label string) error {
	log.Println("we enter create notification")
	if ticket.ID == 0 {
		return errors.New("no ticket found")
	}
	notification := models.Notification{
		Title:    "New Ticket Request",
		Message:  fmt.Sprintf("%s has asked for a new %s for project %s. ", ticket.Email, ticket.Type, ticket.Project.Name),
		New:      true,
		TicketID: ticket.ID,
	}

	if result := ns.DB.Create(&notification); result != nil {
		return result.Error
	}
	if err := manager.SendClientNotification("refetch"); err != nil {
		return err
	}

	return nil
}
