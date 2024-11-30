package services

import (
	"errors"
	"log"
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

func (ss *StatisticsServices) GetYearly(year string) ([]models.Stats, error) {
	log.Println("wtf is this??")
	var statistics []models.Stats
	result := ss.DB.Where("year = ?", year).Find(&statistics)
	if result.Error != nil {
		return []models.Stats{}, result.Error
	}
	return statistics, nil
}

func (ss *StatisticsServices) GetMonthly(month, year string) (models.Stats, error) {
	var statistics models.Stats
	result := ss.DB.Where("year = ?", year).Find(&statistics)
	if result.Error != nil {
		return ss.Statistics, result.Error
	}
	return statistics, nil
}

func (ss *StatisticsServices) HandleStatsSubscription(oldSubscription, subscription models.Subscription) error {
	if oldSubscription.ID == 0 && subscription.PaymentStatus != "NOT_PAID" {
		subscription.ID = 0
		if subscription.PaymentStatus == "TOPAY" {
			subscription.Amount = 0
		}
		err := ss.AddStatsSubscription(subscription)
		if err != nil {
			return err
		}
		return nil
	}

	if oldSubscription.PaymentStatus != "PAID" && subscription.PaymentStatus == "PAID" {
		err := ss.AddStatsSubscription(subscription)
		if err != nil {
			return err
		}
		return nil
	}
	if oldSubscription.PaymentStatus == "PAID" && subscription.PaymentStatus != "PAID" {
		err := ss.RemoveStatsSubscription(subscription)
		if err != nil {
			return err
		}
		return nil
	}
	if oldSubscription.PaymentStatus == "PAID" && subscription.PaymentStatus == "PAID" && oldSubscription.Amount != subscription.Amount {
		subscription.Amount = subscription.Amount - oldSubscription.Amount
		err := ss.AddStatsSubscription(subscription)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (ss *StatisticsServices) AddStatsSubscription(subscription models.Subscription) error {
	today := time.Now()
	currentMonth := int(today.Month())
	currentYear := today.Year()
	var stats models.Stats
	ss.DB.Find(&stats, "month = ? AND year = ?", currentMonth, currentYear)

	if subscription.ID == 0 {
		stats.TotalSubscriptions += 1
	}
	if subscription.Tag == "NEW" {
		stats.NewSubscriptions += 1
	}
	stats.TotalRevenue += subscription.Amount
	stats.NetProfit += subscription.Amount

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

func (ss *StatisticsServices) RemoveStatsSubscription(subscription models.Subscription) error {
	currentMonth := int(subscription.PaidAt.Month())
	currentYear := subscription.PaidAt.Year()
	var stats models.Stats
	ss.DB.Find(&stats, "month = ? AND year = ?", currentMonth, currentYear)
	stats.TotalRevenue -= subscription.Amount
	stats.NetProfit -= subscription.Amount
	if subscription.ID == 0 {
		stats.TotalSubscriptions -= 1
	}

	if subscription.Tag == "NEW" {
		stats.NewSubscriptions -= 1
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

func (ss *StatisticsServices) InvoiceStats(invoice models.Invoice, operation string) error {
	currentMonth := int(invoice.CreatedAt.Month())
	currentYear := invoice.CreatedAt.Year()
	var stats models.Stats
	ss.DB.Find(&stats, "month = ? AND year = ?", currentMonth, currentYear)
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
	ss.DB.Find(&stats, "month = ? AND year = ?", currentMonth, currentYear)
	if invoice.InvoiceType == "PAYOUT" {
		stats.TotalExpenses -= invoice.Amount
	} else {
		stats.TotalRevenue -= invoice.Amount
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

func (ss *StatisticsServices) ScriptStatsAppointment() error {
	var subs []models.Subscription
	ss.DB.Where("payment_status = ?", "PAID").Find(&subs)
	firstRun := true
	firstMonth := int(subs[0].PaidAt.Month())
	for _, app := range subs {
		currentMonth := int(app.PaidAt.Month())
		currentYear := app.PaidAt.Year()
		if app.PaidAt != nil {
			currentMonth = int(app.PaidAt.Month())
			currentYear = app.PaidAt.Year()
		}
		var stats models.Stats
		if currentMonth > firstMonth {
			firstRun = true
			firstMonth = currentMonth
		}
		ss.DB.Find(&stats, "month = ? AND year = ?", currentMonth, currentYear)
		if firstRun {
			stats.TotalRevenue = 0
			stats.TotalSubscriptions = 0
			firstRun = false
		}
		stats.TotalRevenue += app.Amount
		stats.TotalSubscriptions += 1
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
	}

	return nil
}
