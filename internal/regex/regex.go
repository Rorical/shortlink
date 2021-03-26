package regex

import "regexp"

var (
	urlexp, _   = regexp.Compile("(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")
	properid, _ = regexp.Compile("^[_a-zA-Z0-9]+$")
)

func IsUrl(str string) bool {
	return urlexp.MatchString(str)
}

func IsProperId(str string) bool {
	return properid.MatchString(str)
}
