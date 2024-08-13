package pkg

import "regexp"

func IsValidGUID(str string) bool {

	pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}
