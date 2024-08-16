package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	Name            string        `json:"name" gorm:"unique"`
	File            string        `json:"file"`
	TotalUser       float64       `json:"total_user"`
	TotalActiveUser string        `json:"total_active_user"`
	Notes           string        `json:"notes"`
	DomainURL       string        `json:"domain_url"`
	SSLExpiredAt    time.Time     `json:"ssl_expired_at"`
	DBName          string        `json:"db_name"`
	RepoName        string        `json:"repo_name"`
	CommitStats     []CommitStats `json:"commit_stats" gorm:"foreignKey:ProjectID"`
	Status          string        `json:"status" gorm:"default:ACTIVE"`
	UpdateCommands  string        `json:"commands" gorm:"type:text"`
	BackupAt        *time.Time    `json:"backup_at"`
	LastBuildAt     *time.Time    `json:"last_build_at"`
	CreatedAt       time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
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
