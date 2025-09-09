package rules

import (
	"fmt"
	"log"
	"strings"

	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/contracts/validation"
)

type Filetype struct{}

func (r *Filetype) Signature() string {
	return "filetype"
}

func (r *Filetype) Passes(data validation.Data, val any, options ...any) bool {
	// val adalah *filesystem.File
	filePtr, ok := val.(**filesystem.File) // pointer ke interface
	if !ok || filePtr == nil {
		return false
	}

	file := *filePtr // dereference pointer → interface
	ext, err := (*file).Extension()
	if err != nil {
		return false
	}

	ext = strings.ToLower(ext)
	log.Println("ekstensi file:", ext)

	if len(options) == 0 {
		return false
	}

	for _, opt := range options {
		if strings.ToLower(fmt.Sprint(opt)) == ext {
			return true
		}
	}

	return false
}

func (r *Filetype) Message() string {
	return "file harus bertipe sesuai ekstensi yang diperbolehkan"
}
