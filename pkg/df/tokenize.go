package df

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Tokenize splits s into tokens separated by spaces. Spaces in nested quotes
// strings are not split, thus 'Java Dev Hamburg' becomes one token.
func Tokenize(s string) []string {
	var tokens []string
	t := ""
	quoted := false
	for i := 0; i < len(s); i++ {
		if string(s[i]) == "'" {
			quoted = !quoted
			continue
		}

		if string(s[i]) != " " {
			t = fmt.Sprintf("%s%s", t, string(s[i]))
			continue
		}

		if string(s[i]) == " " && quoted {
			t = fmt.Sprintf("%s%s", t, string(s[i]))
			continue
		}

		if string(s[i]) == " " {
			tokens = append(tokens, t)
			t = ""
		}
	}
	tokens = append(tokens, t)
	log.Debug(tokens)
	return tokens
}
