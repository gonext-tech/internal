package services

import (
	"errors"
	"time"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type StatisticsServices struct {
	Statistics models.Stats
	DB         *gorm.DB
}

func NewStatisticcServices(s models.Stats, db *gorm.DB) *StatisticsServices {
	return &StatisticsServices{
		Statistics: s,
		DB:         db,
	}
}

func (ss *StatisticsServices) GetALL(month, year string, user models.User) (models.Stats, error) {
	var statistics models.Stats
	result := ss.DB.Where("shop_id = ? AND month = ? AND year = ?", user.ShopID, month, year).Find(&statistics)
	if result.Error != nil {
		return models.Stats{}, result.Error
	}
	return statistics, nil
}

func (ss *StatisticsServices) GetDaily(startDate, endDate string, user models.User) (models.Stats, error) {
	var subscriptions []models.Subscription
	var invoices []models.Invoice

	query := ss.DB
	if startDate != "" && endDate != "" {
		var fromDate, toDate time.Time
		layout := "2006-01-02"

		fromDate, err := time.Parse(layout, startDate)
		if err != nil {
			panic(err)
		}
		toDate, err = time.Parse(layout, endDate)
		if err != nil {
			panic(err)
		}
		query = query.Where("date BETWEEN ? AND ?", fromDate, toDate.AddDate(0, 0, 1))
	}
	result := query.Where("shop_id = ? AND payment_status = ?", user.ShopID, "PAID").Find(&subscriptions)
	query.Where("shop_id = ? ", user.ShopID).Find(&invoices)
	var totalSubscriptions int
	var totalRevenue, totalExpenses float64
	for _, subscriptions := range subscriptions {
		totalRevenue += subscriptions.Amount
	}
	for _, invoice := range invoices {
		if invoice.InvoiceType == "PAYOUT" {
			totalExpenses += invoice.Amount
		} else {
			totalRevenue += invoice.Amount
		}
	}
	totalSubscriptions = len(subscriptions)
	if result.Error != nil {
		return models.Stats{}, result.Error
	}

	statistics := models.Stats{
		TotalRevenue:       totalRevenue,
		TotalExpenses:      totalExpenses,
		TotalSubscriptions: totalSubscriptions,
		NetProfit:          totalRevenue - totalExpenses,
	}
	return statistics, nil
}

func (ss *StatisticsServices) GetShop(shopID string) ([]models.Stats, error) {
	var statistics []models.Stats
	result := ss.DB.Where("shop_id = ?", shopID).Find(&statistics)
	if result.Error != nil {
		return []models.Stats{}, result.Error
	}
	return statistics, nil
}

func (ss *StatisticsServices) AddStatsSubscription(subscription models.Subscription) error {
	today := time.Now()
	currentMonth := int(today.Month())
	currentYear := today.Year()
	var stats models.Stats
	ss.DB.Find(&stats, "shop_id = ? AND month = ? AND year = ?", subscription.ShopID, currentMonth, currentYear)
	stats.TotalRevenue += subscription.Amount
	stats.TotalSubscriptions += 1
	if stats.ID == 0 {
		stats.Month = currentMonth
		stats.ShopID = subscription.ShopID
		stats.Year = currentYear
		if result := ss.DB.Create(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	} else {
		if result := ss.DB.Save(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	}
	return nil
}

func (ss *StatisticsServices) RemoveStatsAppointment(subscription models.Subscription) error {
	currentMonth := int(subscription.PaidAt.Month())
	currentYear := subscription.PaidAt.Year()
	var stats models.Stats
	ss.DB.Find(&stats, "shop_id = ? AND month = ? AND year = ?", subscription.ShopID, currentMonth, currentYear)
	stats.TotalRevenue -= subscription.Amount
	stats.NetProfit -= subscription.Amount
	stats.TotalSubscriptions -= 1

	if stats.ID == 0 {
		stats.Month = currentMonth
		stats.ShopID = subscription.ShopID
		stats.Year = currentYear
		if result := ss.DB.Create(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	} else {
		if result := ss.DB.Save(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	}
	return nil
}

func (ss *StatisticsServices) InvoiceStats(invoice models.Invoice, operation string) error {
	currentMonth := int(invoice.CreatedAt.Month())
	currentYear := invoice.CreatedAt.Year()
	var stats models.Stats
	ss.DB.Find(&stats, "shop_id = ? AND month = ? AND year = ?", invoice.ShopID, currentMonth, currentYear)
	stats.ShopID = invoice.ShopID
	if operation == "ADD" {
		if invoice.InvoiceType == "PAYOUT" {
			stats.TotalExpenses += invoice.Amount
		} else {
			stats.TotalRevenue += invoice.Amount
		}
	} else {
		if invoice.InvoiceType == "PAYOUT" {
			stats.TotalExpenses -= invoice.Amount
		} else {
			stats.TotalRevenue -= invoice.Amount
		}
	}
	if stats.ID == 0 {
		stats.Month = currentMonth
		stats.Year = currentYear
		if result := ss.DB.Create(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	} else {
		if result := ss.DB.Save(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	}

	return nil
}

func (ss *StatisticsServices) RemoveStatsExpense(invoice models.Invoice) error {
	currentMonth := int(invoice.CreatedAt.Month())
	currentYear := invoice.CreatedAt.Year()
	var stats models.Stats
	ss.DB.Find(&stats, "shop_id = ? AND month = ? AND year = ?", invoice.ShopID, currentMonth, currentYear)
	if invoice.InvoiceType == "PAYOUT" {
		stats.TotalExpenses -= invoice.Amount
	} else {
		stats.TotalRevenue -= invoice.Amount
	}
	if stats.ID == 0 {
		stats.Month = currentMonth
		stats.ShopID = invoice.ShopID
		stats.Year = currentYear
		if result := ss.DB.Create(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	} else {
		if result := ss.DB.Save(&stats); result.Error != nil {
			return errors.New(result.Error.Error())
		}
	}
	return nil
}

func (ss *StatisticsServices) ScriptStatsAppointment() error {
	var apts []models.Appointment
	ss.DB.Where("payment_status = ?", "PAID").Find(&apts)
	firstRun := true
	firstMonth := int(apts[0].Date.Month())
	for _, app := range apts {
		currentMonth := int(app.Date.Month())
		currentYear := app.Date.Year()
		if app.PaidAt != nil {
			currentMonth = int(app.PaidAt.Month())
			currentYear = app.PaidAt.Year()
		}
		var stats models.Stats
		if currentMonth > firstMonth {
			firstRun = true
			firstMonth = currentMonth
		}
		ss.DB.Find(&stats, "shop_id = ? AND month = ? AND year = ?", app.ShopID, currentMonth, currentYear)
		if firstRun {
			stats.TotalRevenue = 0
			stats.TotalSubscriptions = 0
			firstRun = false
		}
		stats.TotalRevenue += app.Price
		stats.TotalSubscriptions += 1
		if stats.ID == 0 {
			stats.Month = currentMonth
			stats.ShopID = app.ShopID
			stats.Year = currentYear
			if result := ss.DB.Create(&stats); result.Error != nil {
				return errors.New(result.Error.Error())
			}
		} else {
			if result := ss.DB.Save(&stats); result.Error != nil {
				return errors.New(result.Error.Error())
			}
		}
	}

	return nil
}
