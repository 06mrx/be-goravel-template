package models

import (
	"goravel_api/app/auditable"
	"time"

	"github.com/google/uuid"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
	// "github.com/goravel/framework/contracts/auth"
)

type User struct {
	orm.Model
	ID           uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name         string     `gorm:"column:name"`
	Email        string     `gorm:"column:email;unique"`
	Password     string     `gorm:"column:password" json:"-"`
	GoogleId     string     `gorm:"column:google_id"`
	AuthProvider string     `gorm:"column:auth_provider"`
	CreatedBy    string     `gorm:"column:created_by"`
	UpdatedBy    string     `gorm:"column:updated_by"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	Roles        []Role     `gorm:"many2many:role_user;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:role_id"`
	Articles     []Article  `gorm:"foreignKey:CreatedBy;references:ID" json:"articles"`
	orm.SoftDeletes
	orm.Timestamps
	auditable.Auditable
}

// Metode pembantu untuk memeriksa izin
func (u *User) HasPermissionTo(permissionName string) bool {
	// Muat relasi peran
	facades.Orm().Query().With("Roles.Permissions").Find(u)

	// Iterasi peran dan izin untuk menemukan kecocokan
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true
			}
		}
	}
	return false
}
