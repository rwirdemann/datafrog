package matcher

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

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
