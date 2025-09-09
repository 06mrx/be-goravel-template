// app/models/audit.go
package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
)

type Audit struct {
	ID            uint   `gorm:"primaryKey"`
	AuditableType string `gorm:"auditable_type"`
	AuditableID   string `gorm:"auditable_id"`
	Event         string
	UserID        string
	OldValues     []byte     `gorm:"type:json"`
	NewValues     []byte     `gorm:"type:json"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
	orm.Timestamps
}

// Fungsi untuk menyimpan audit
func (a *Audit) Save() error {
	return facades.Orm().Query().Create(a)
}
