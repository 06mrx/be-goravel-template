package rules

import (
	"fmt"

	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type Unique struct {
	table  string
	column string
	id     string // Tambahkan field untuk menyimpan ID
}

// Signature Get the name of the rule.
func (r *Unique) Signature() string {
	return "unique"
}

// Passes checks if the value is unique in database
// options[0] = table, options[1] = column (optional), options[2] = id (optional)
func (r *Unique) Passes(data validation.Data, val any, options ...any) bool {
	// fmt.Printf(""") // Debug output
	if len(options) == 0 {
		return false // Harus ada table
	}

	r.table, _ = options[0].(string)
	if len(options) > 1 {
		r.column, _ = options[1].(string)
	} else {
		r.column = r.table // Fallback jika kolom tidak ditentukan
	}

	// 1. Dapatkan ID yang akan diabaikan dari opsi
	if len(options) > 2 {
		r.id, _ = options[2].(string)
	}
	fmt.Printf("ID to ignore: %s\n", r.id) // Debug output
	// fmt.Printf("Table: %s, Column: %s, ID to ignore: %s, Value: %v\n", r.table, r.column, r.id, val) // Debug output
	// 2. Buat kueri database
	query := facades.DB().Table(r.table).Where(r.column, val)

	// 3. Tambahkan kondisi `deleted_at` secara kondisional
	// Cek keberadaan kolom deleted_at
	var hasSoftDelete bool
	hasSoftDelete = facades.Schema().HasColumn(r.table, "deleted_at")

	if hasSoftDelete {
		query = query.WhereNull("deleted_at")
	}

	if hasSoftDelete {
		query = query.WhereNull("deleted_at")
	}
	if r.id != "" {
		query = query.Where("id <> ?", r.id)
	}

	// 4. Hitung hasil
	var count int64
	count, err := query.Count()
	// fmt.Printf("Count: %d\n", count)
	fmt.Printf("Error: %v\n", err) // Debug output
	if err != nil {
		return false
	}
	fmt.Printf("Count: %d\n", count) // Debug output

	return count == 0
}

// Message Get the validation error message.
func (r *Unique) Message() string {
	// return fmt.Sprintf("The value must be unique in table '%s', column '%s'.", r.table, r.column)
	return "The :attribute has already been taken."
}
