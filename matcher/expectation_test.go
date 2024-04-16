package matcher

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		desc  string
		s     string
		count int
	}{
		{
			desc:  "",
			s:     "select * from jobs where id=1;",
			count: 6,
		},
		{
			desc:  "",
			s:     "select job0_.id as id1_0_, job0_.description as descript2_0_ from job job0_ order by job0_.publish_at desc",
			count: 14,
		},
		{
			desc:  "",
			s:     "update job set description='World, X' where id=39",
			count: 6,
		},
		{
			desc:  "",
			s:     "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World, Huhu', '2024-04-02 08:37:37', 0, null, '', 'Hello', 39)",
			count: 18,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tokens := Tokenize(tC.s)
			for _, v := range tokens {
				log.Printf("%v", v)
			}
			assert.Len(t, tokens, tC.count)
		})
	}
}

func TestBuildDiff(t *testing.T) {
	tokens := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:48:15', 0, null, '', 'Hello', 2)"
	expectation := Expectation{Tokens: Tokenize(tokens)}

	v := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:48:16', 0, null, '', 'Hello', 3)"
	diff, _ := expectation.BuildDiff(v)
	assert.Len(t, diff, 2)
	assert.Equal(t, 12, diff[0])
	assert.Equal(t, 17, diff[1])
}
