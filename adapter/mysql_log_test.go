package adapter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	s := "2024-04-02T06:38:05.015501Z	1669 Query	update job set description='World, X', publish_at='2024-04-02 08:37:37', tags='', title='Hello' where id=39"
	actual, err := MySQLLog{}.Timestamp(s)
	assert.Nil(t, err)
	expected, err := time.Parse(time.RFC3339Nano, "2024-04-02T06:38:05.015501Z")
	assert.Equal(t, expected, actual)
}
