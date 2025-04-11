package models

import "time"

type Domain struct {
	ID             uint            `form:"id" gorm:"primaryKey"`
	Name           string          `gorm:"unique" form:"name"`
	Provider       string          `form:"provider"`
	AccountEmail   string          `form:"account_email"`
	ServerID       uint            `form:"server_id"`
	Server         MonitoredServer `gorm:"foreignKey:ServerID" form:"server"`
	ExpirationDate time.Time       `form:"expiration_date"`
	RenewalDate    time.Time       `form:"renewal_date"`
	AutoRenew      bool            `form:"auto_renew"`
	SSLExpiryDate  *time.Time      `form:"ssl_expiry_date"`
	Status         string          `form:"status"`
	AnnualCost     float64         `form:"annual_cost"`
	Notes          string          `gorm:"type:text" form:"notes"`
}
