package models

import (
	"gorm.io/gorm"
	"time"
)

type ServerStatus string

const (
	StatusUp   ServerStatus = "up"
	StatusDown ServerStatus = "down"
)

type MonitoredServer struct {
	gorm.Model

	// Core info
	Name      string `gorm:"not null" form:"name"`
	IPAddress string `gorm:"uniqueIndex;not null" form:"ip_address"`
	Hostname  string `form:"hostname"`
	Location  string `form:"location"`

	// Ownership
	SignupEmail string    `gorm:"not null" form:"signup_email"`
	RenewalDate time.Time `form:"renewal_date" `
	AnnualCost  float64   `form:"annual_cost"`

	// Hardware specs
	RAMGB     int `form:"ram_gb"`
	StorageGB int `form:"storage_gb"`

	// Monitoring data
	Status     ServerStatus   `gorm:"type:varchar(10);default:'down'" form:"status"`
	LastSeen   *time.Time     `form:"last_seen"`
	Uptime     *time.Duration `form:"uptime"`
	Downtime   *time.Duration `form:"downtime"`
	LastChange *time.Time     `form:"last_change"`
	Notes      string         `form:"notes"`

	// Logs
	//Logs []ServerLog `gorm:"foreignKey:ServerID" json:"logs"`
}

type ServerLog struct {
	gorm.Model
	ServerID  uint         `gorm:"index" json:"server_id"`
	Timestamp time.Time    `gorm:"autoCreateTime" json:"timestamp"`
	Status    ServerStatus `gorm:"type:varchar(10)" json:"status"`
	Message   string       `json:"message"`
}
