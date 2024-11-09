package utils

import "regexp"

func IsValidUrl(url string) bool {
	pattern := `^((http|https):\/\/)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(:[0-9]{1,5})?(\/[^\s]*)?$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(url)
}

func IsValidFilePath(path string) bool {
	pattern := `^([a-zA-Z]:\\|\/)?([a-zA-Z0-9._-]+[\/\\]?)*$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(path)
}

func OneOfThem[T comparable](a T, b []T) bool {
	for i := range b {
		if a == b[i] {
			return true
		}
	}

	return false
}

func ConvertNumericUnits(numericUnit string) string {
	if numericUnit == "millisecond" {
		return "ms"
	} else if numericUnit == "second" {
		return "s"
	} else if numericUnit == "minute" {
		return "m"
	} else if numericUnit == "hour" {
		return "h"
	} else {
		return ""
	}
}
