package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	// s contains tabs
	s := "2024-04-02T06:38:05.015501Z	1669 Query	update job set description='World, X', publish_at='2024-04-02 08:37:37', tags='', title='Hello' where id=39"
	_, err := Timestamp(s)
	assert.Nil(t, err)
}
