package learn_lib

import (
	"regexp"
)

func IsEnglishWord(word string) bool {
	st, err := regexp.Match("^\\w+$", []byte(word))
	if err != nil {
		return false
	}
	return st

}
