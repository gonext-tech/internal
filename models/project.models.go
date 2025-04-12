package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID              uint            `form:"id" gorm:"primaryKey"`
	Name            string          `form:"name" gorm:"unique"`
	File            string          `form:"file"`
	TotalUser       float64         `form:"total_user"`
	TotalActiveUser string          `form:"total_active_user"`
	Notes           string          `form:"notes"`
	DomainURL       string          `form:"domain_url"`
	DBName          string          `form:"db_name"`
	RepoName        string          `form:"repo_name"`
	CommitStats     []CommitStats   `form:"-" gorm:"foreignKey:ProjectID"`
	ServerID        uint            `form:"server_id"`
	Server          MonitoredServer `form:"-" gorm:"foreignKey:ServerID"`
	LeadID          uint            `form:"lead_id"`
	Lead            Admin           `form:"-" gorm:"foreignKey:LeadID"`
	Status          string          `form:"status" gorm:"default:ACTIVE"`
	UpdateCommands  string          `form:"commands" gorm:"type:text"`
	BackupAt        *time.Time      `form:"backup_at"`
	LastBuildAt     *time.Time      `form:"last_build_at"`
	CreatedAt       time.Time       `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time       `form:"updated_at" gorm:"autoUpdateTime"`
}

type GithubStats struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	ProjectID uint `json:"project_id"`
}

type Folder struct {
	Path     string    `json:"path"`
	Name     string    `json:"name,omitempty"`
	IsDir    bool      `json:"isDir"`
	Size     int64     `json:"size,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	IsEmpty  bool      `json:"isEmpty,omitempty"`
	Children []Folder  `json:"folders,omitempty"`
}

type ProjectsDB struct {
	Name string
	DB   *gorm.DB
}

type CommitStats struct {
	ID           uint       `json:"id"`
	ProjectID    uint       `json:"project_id"`
	SHA          *string    `json:"sha,omitempty"`
	AuthorName   *string    `json:"author_name,omitempty"`
	AuthorEmail  *string    `json:"author_email,omitempty"`
	AuthorAvatar *string    `json:"author_avatar,omitempty"`
	Message      *string    `json:"message,omitempty"`
	Additions    *int       `json:"additions,omitempty"`
	Deletions    *int       `json:"deletions,omitempty"`
	Compare      *string    `json:"compare"`
	Modified     *int       `json:"modified,omitempty"`
	Total        *int       `json:"total,omitempty"`
	HTMLURL      *string    `json:"html_url,omitempty"`
	URL          *string    `json:"url,omitempty"`
	NodeID       *string    `json:"node_id,omitempty"`
	CommentCount *int       `json:"comment_count,omitempty"`
	Timestamp    *time.Time `json:"timestamp,omitempty"`
}
