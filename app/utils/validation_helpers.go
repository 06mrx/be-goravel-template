package utils

import "github.com/goravel/framework/contracts/filesystem"

// func CheckFileType (fileHeader *multipart.FileHeader, allowedTypes []string) bool {

// 	for _, ext := range allowedTypes {
// 		if strings.ToLower(extension) == strings.ToLower(ext) {
// 			return true
// 		}
// 	}
// }

func CheckFileTypeAndSize(allowedExt []string, maxUploadSize int64, file *filesystem.File) bool {
	extension, _ := (*file).Extension()
	size, _ := (*file).Size()
	extCheck := false
	sizeCheck := false
	for _, a := range allowedExt {
		if a == extension {
			extCheck = true
		}
	}
	if size <= maxUploadSize {
		sizeCheck = true
	}
	if extCheck && sizeCheck {
		return true
	}
	return false
}
