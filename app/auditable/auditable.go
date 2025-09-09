// app/auditable/auditable.go
package auditable

import (
	"context"
	"encoding/json"
	"fmt"

	// "fmt"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditData adalah struct lokal untuk menyimpan data audit.
// Ini bukan model GORM, hanya struktur data.
type AuditData struct {
	AuditableType string `json:"auditable_type"`
	AuditableID   string `json:"auditable_id"`
	Event         string `json:"event"`
	UserID        string `json:"user_id"`
	OldValues     []byte `json:"old_values"`
	NewValues     []byte `json:"new_values"`
}

// Auditable adalah trait yang dapat disematkan ke model.
type Auditable struct{}

// AfterCreate dipanggil setelah record dibuat.
func (a *Auditable) AfterCreate(tx *gorm.DB) error {
	if tx.Error != nil {
		return tx.Error
	}

	data := GetFieldsFromModel(tx.Statement.Model, "ID", "CreatedBy")
	// Dapatkan data baru dari model
	newData, _ := json.Marshal(tx.Statement.Model)

	audit := &AuditData{
		AuditableType: tx.Statement.Schema.Table,
		AuditableID:   data["ID"].(uuid.UUID).String(),
		Event:         "created",
		UserID:        data["CreatedBy"].(string),
		NewValues:     newData,
	}

	// Simpan data audit langsung ke tabel
	return tx.Table("audits").Create(audit).Error
}

// 💡 Callback BeforeUpdate: Ambil data lama sebelum update
func (a *Auditable) BeforeUpdate(tx *gorm.DB) error {
	// Dapatkan nilai primary key dari record yang akan diupdate
	data := GetFieldsFromModel(tx.Statement.Model, "ID", "UpdatedBy")
	id := data["ID"]

	// Gunakan GORM untuk mencari record lama berdasarkan ID-nya
	// Pastikan model yang digunakan adalah model yang sama dengan yang sedang diupdate
	model := reflect.New(reflect.TypeOf(tx.Statement.Model).Elem()).Interface()

	// 💡 Lakukan kueri pada transaksi yang sama (`tx`).
	// Jangan membuat sesi baru.
	if err := tx.Table(tx.Statement.Schema.Table).Where("id", id).First(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	newContext := context.WithValue(tx.Statement.Context, "old_values", model)
	tx.Statement.Context = newContext
	return nil
}

// AfterUpdate dipanggil setelah record diperbarui.
func (a *Auditable) AfterUpdate(tx *gorm.DB) error {
	if tx.Error != nil {
		return tx.Error
	}

	newData, _ := json.Marshal(tx.Statement.Model)
	data := GetFieldsFromModel(tx.Statement.Model, "ID", "UpdatedBy")
	// Untuk mendapatkan data lama, kita perlu melakukan kueri tambahan di BeforeUpdate.
	// Jika tidak ada data lama yang disimpan, OldValues akan kosong.
	// 💡 Ambil data dari context.Context
	var oldData []byte
	if old, ok := tx.Statement.Context.Value("old_values").(interface{}); ok && old != nil {
		oldData, _ = json.Marshal(old)
	}

	fmt.Println(oldData)

	audit := &AuditData{
		AuditableType: tx.Statement.Schema.Table,
		AuditableID:   data["ID"].(uuid.UUID).String(),
		Event:         "updated",
		UserID:        data["UpdatedBy"].(string),
		OldValues:     oldData,
		NewValues:     newData,
	}

	return tx.Table("audits").Create(audit).Error
	// return tx.Error
}

// AfterDelete dipanggil setelah record dihapus.
func (a *Auditable) AfterDelete(tx *gorm.DB) error {
	if tx.Error != nil {
		return tx.Error
	}

	// Dapatkan data lama (yang baru saja dihapus)
	data := GetFieldsFromModel(tx.Statement.Model, "ID", "UpdatedBy")
	oldData, _ := json.Marshal(tx.Statement.Model)

	audit := &AuditData{
		AuditableType: tx.Statement.Schema.Table,
		AuditableID:   data["ID"].(uuid.UUID).String(),
		Event:         "deleted",
		UserID:        data["UpdatedBy"].(string),
		OldValues:     oldData,
	}
	return tx.Table("audits").Create(audit).Error
}

func GetFieldsFromModel(model any, fields ...string) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(model)

	// Jika pointer, ambil elemennya
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, field := range fields {
		f := val.FieldByName(field)
		if f.IsValid() {
			result[field] = f.Interface()
		}
	}
	return result
}
