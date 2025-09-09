package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/goravel/framework/facades"
	"gorm.io/gorm"
)

// find model by id
// "name ILIKE ?", "%"+searchQuery+"%"
// func FindModelByID(id string, model any) error {
// 	err := facades.Orm().Query().Where("id", id).First(model)
// 	if errors.Is(err, gorm.ErrRecordNotFound) || err != nil {
// 		return err
// 	}

//		fmt.Printf("Found model by model helper: %+v\n", model)
//		return nil
//	}
// func FindModelByID(id string, model any) error {
// 	err := facades.Orm().Query().Where("id", id).First(model)
// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return err
// 	}

// 	// cek apakah field ID = UUID kosong (0000...)
// 	// gunakan reflect biar bisa ke semua model
// 	val := reflect.ValueOf(model).Elem()
// 	idField := val.FieldByName("ID")
// 	if idField.IsValid() {
// 		if idStr, ok := idField.Interface().(uuid.UUID); ok {
// 			if idStr == uuid.Nil {
// 				return fmt.Errorf("data dengan id %s tidak ditemukan", id)
// 			}
// 		}
// 	}

// 	// fmt.Printf("Found model by model helper: %+v\n", model)
// 	return nil
// }

func FindModelByID(id string, model any, relations ...string) error {
	query := facades.Orm().Query().Where("id", id)

	// Parameter variadik 'relations' akan otomatis menjadi slice.
	// GORM akan memuat kolom tertentu dari relasi jika sintaksnya seperti "Roles:id".
	for _, relation := range relations {
		query = query.With(relation)
	}

	err := query.First(model)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	val := reflect.ValueOf(model).Elem()
	idField := val.FieldByName("ID")
	if idField.IsValid() {
		if idStr, ok := idField.Interface().(uuid.UUID); ok {
			if idStr == uuid.Nil {
				return fmt.Errorf("data dengan id %s tidak ditemukan", id)
			}
		}
	}

	return nil
}
