package matcher

import (
	"testing"

	"github.com/rwirdemann/databasedragon/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestNewTokenMatcher(t *testing.T) {
	c := config.NewConfig("config.json")
	expectations := []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;"}
	verifications := []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=2;"}
	tm := NewTokenMatcher(c, expectations, verifications)
	assert.Equal(t, 1, len(tm.expectations))
	e := tm.expectations[0]
	assert.Equal(t, 6, len(e.Tokens))
}

func TestMatches(t *testing.T) {
	testCases := []struct {
		desc          string
		expectations  []string
		verifications []string
		actual        string
		index         int
	}{
		{
			desc:          "",
			expectations:  []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;"},
			verifications: []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=2;"},
			actual:        "2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;",
			index:         0,
		},
		{
			desc:          "",
			expectations:  []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;"},
			verifications: []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=2;"},
			actual:        "2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=3;",
			index:         0,
		},
		{
			desc:          "",
			expectations:  []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;"},
			verifications: []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=2;"},
			actual:        "2024-04-02T06:37:24.479429Z 1669 Query update * from job where id=3;",
			index:         -1,
		},
		{
			desc: "",
			expectations: []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;",
				"2024-04-02T06:38:05.015501Z 1669 Query	update job set description='World, X', publish_at='2024-04-02 08:37:37', tags='', title='Hello' where id=39"},
			verifications: []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=2;",
				"2024-04-02T06:39:05.015501Z 1669 Query	update job set description='World, X', publish_at='2024-04-02 08:37:39', tags='', title='Hello' where id=40"},
			actual: "2024-04-02T06:40:05.015501Z 1669 Query	update job set description='World, X', publish_at='2024-04-02 08:37:40', tags='', title='Hello' where id=41",
			index:  1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c := config.NewConfig("config.json")
			tm := NewTokenMatcher(c, tC.expectations, tC.verifications)
			assert.Equal(t, tC.index, tm.Matches(tC.actual))
		})
	}
}

func TestRemoveExpectation(t *testing.T) {
	expectations := []string{"select * from job where id=1;", "update job set description='World, X'where id=39"}
	verifications := []string{"select * from job where id=12;", "update job set description='World, X'where id=40"}
	tm := NewTokenMatcher(config.NewConfig("config.json"), expectations, verifications)
	tm.RemoveExpectation(0)
	assert.Len(t, tm.expectations, 1)
	tm.RemoveExpectation(0)
	assert.Len(t, tm.expectations, 0)
}
