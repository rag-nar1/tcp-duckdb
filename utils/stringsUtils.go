package utils

import "strings"


func Trim(s string) string {
	for len(s) > 0 && !(s[0] >= 33 && s[0] <= 126) {
		s = s[1:]
	}
	for len(s) > 0 && !(s[len(s) - 1] >= 33 && s[len(s) - 1] <= 126) {
		s = s[:len(s) - 1]
	}
	return s
}

func ToLower(s ...string) []string {
	for i := range s {
		s[i] = strings.ToLower(s[i])
	}

	return s
}

func TrimList(list []string) {
	for i := 0; i < len(list); i ++ {
		list[i] = Trim(list[i])
	}
}	