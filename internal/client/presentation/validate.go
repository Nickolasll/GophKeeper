package presentation

import "regexp"

var validCardNumber, validValidThru, validCVV, validCardHolder *regexp.Regexp

func regexpMustCompile() {
	validCardNumber = regexp.MustCompile(`\b(\d{4}\s\d{4}\s\d{4}\s\d{4}$)\b`)
	validValidThru = regexp.MustCompile(`(0[1-9]|1[012])/\d{2}`)
	validCVV = regexp.MustCompile(`^\d{3,4}$`)
	validCardHolder = regexp.MustCompile(`^((?:[A-Za-z]+ ?){0,3})$`)
}
