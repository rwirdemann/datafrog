package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDiff(t *testing.T) {
	tokens := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:48:15', 0, null, '', 'Hello', 2)"
	expectation := Expectation{Tokens: Tokenize(tokens)}

	v := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:48:16', 0, null, '', 'Hello', 3)"
	diff, _ := expectation.Diff(Tokenize(v))
	assert.Len(t, diff, 2)
	assert.Equal(t, 12, diff[0])
	assert.Equal(t, 17, diff[1])
}
