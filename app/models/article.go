package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/goravel/framework/database/orm"
)

type Article struct {
	orm.Model
	ID          uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title       string     `gorm:"column:title" json:"title"`
	Slug        string     `gorm:"column:slug;unique" json:"slug"`
	ImageUrl    string     `gorm:"column:image_url" json:"image_url"`
	Content     string     `gorm:"column:content" json:"content"`
	Status      string     `gorm:"column:status" json:"status"`
	UserID      string     `gorm:"column:user_id" json:"user_id"`
	PublishedAt *time.Time `gorm:"column:published_at" json:"published_at"`
	CreatedBy   string     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   string     `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	User        *User      `gorm:"foreignKey:UserID;references:ID" json:"user"`

	orm.SoftDeletes
	orm.Timestamps
	// auditable.Auditable
}
