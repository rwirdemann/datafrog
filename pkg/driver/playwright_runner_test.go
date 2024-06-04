package driver

import (
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToPlaywright(t *testing.T) {
	r := NewPlaywrightRunner(df.Config{})
	assert.Equal(t, "full-12.spec.ts", r.ToPlaywright("full-12.json"))
}
