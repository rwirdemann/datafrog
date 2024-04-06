package matcher

import (
	"testing"

	"github.com/rwirdemann/databasedragon/config"
	"github.com/stretchr/testify/assert"
)

func TestNewTokenMatcher(t *testing.T) {
	c := config.NewConfig("config.json")
	expectations := []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;"}
	verifications := []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=2;"}
	tm := NewTokenMatcher(c, expectations, verifications)
	assert.Equal(t, 1, len(tm.expectations))
	e := tm.expectations[0]
	assert.Equal(t, 6, len(e.tokens))
}

func TestMatches(t *testing.T) {
	c := config.NewConfig("config.json")
	expectations := []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;"}
	verifications := []string{"2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=2;"}
	tm := NewTokenMatcher(c, expectations, verifications)

	assert.True(t, tm.Matches("2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=1;"))
	assert.True(t, tm.Matches("2024-04-02T06:37:24.479429Z 1669 Query select * from job where id=3;"))
	assert.False(t, tm.Matches("2024-04-02T06:37:24.479429Z 1669 Query update * from job where id=3;"))
}
