package validation

import (
	"regexp"
)

func IsEmailValid(email string) bool {
	match, _ := regexp.MatchString("^[\\w-.]+@([\\w-]+\\.)+[\\w-]{2,4}", email)
	return match
}
