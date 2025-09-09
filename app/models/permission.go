package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/goravel/framework/database/orm"
)

type Permission struct {
	ID        uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string     `gorm:"column:name;unique" json:"name"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedBy string     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy string     `gorm:"column:updated_by" json:"updated_by"`
	orm.SoftDeletes

	// auditable.Auditable
}
