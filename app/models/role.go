// app/models/role.go
package models

import (
	"time"

	// Sesuaikan dengan path auditable Anda

	"github.com/google/uuid" // <-- Tambahkan import ini
	// <-- Tambahkan import ini jika belum ada
)

type Role struct {
	ID          uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"` // <-- Ubah tipe menjadi uuid.UUID
	Name        string       `gorm:"unique"`
	Users       []User       `gorm:"many2many:role_user;foreignKey:id;joinForeignKey:role_id;References:id;joinReferences:user_id"`
	Permissions []Permission `gorm:"many2many:permission_role;foreignKey:id;joinForeignKey:role_id;References:id;joinReferences:permission_id"` // Gunakan 'id' bukan 'ID'
	CreatedAt   time.Time    `gorm:"column:created_at"`
	UpdatedAt   time.Time    `gorm:"column:updated_at"`
	// Jika Anda menggunakan gorm.Model atau auditable.Auditable, pastikan ID-nya disematkan atau didefinisikan dengan benar
	// GORM akan menangani ID sebagai primary key secara default jika dinamakan 'ID'
	// auditable.Auditable // Jika Anda menggunakan Auditable, pastikan ID tidak ganda
}

// Catatan: Jika Auditable atau gorm.Model sudah menyediakan ID, Anda mungkin tidak perlu mendeklarasikannya secara eksplisit lagi.
// Namun, jika Anda ingin menggunakan uuid.UUID sebagai ID, deklarasi eksplisit seperti di atas adalah cara yang benar.
