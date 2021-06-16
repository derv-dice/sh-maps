package utils

import "regexp"

// GetRxParams - Get all regexp params from string with provided regular expression
func GetRxParams(rx *regexp.Regexp, str string) (pm map[string]string) {
	if !rx.MatchString(str) {
		return nil
	}
	p := rx.FindStringSubmatch(str)
	n := rx.SubexpNames()
	pm = map[string]string{}
	for i := range n {
		if i == 0 {
			continue
		}

		if n[i] != "" && p[i] != "" {
			pm[n[i]] = p[i]
		}
	}
	return
}
