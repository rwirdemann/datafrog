package postgres

import (
	"fmt"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestRegex(t *testing.T) {
	r := regexp.MustCompile(`\$\d\s=\s'(?:[^']|'')*'|\$\d\s=\sNULL`)
	matches := r.FindAllStringSubmatch("DETAIL:  parameters: $1 = 'World', $2 = '2024-04-19 10:12:12', $3 = '0', $4 = NULL, $5 = '', $6 = 'Hello', $7 = '1'", -1)
	for _, v := range matches {
		fmt.Println(v)
	}
}

func TestReadLine(t *testing.T) {
	c := df.Config{}
	c.Patterns = []string{"insert"}
	pl := NewPostgresLog("postgres.log", c)
	defer pl.Close()
	actual := readLine(t, pl)
	expected := "2024-04-19 10:12:16.889 CEST [89718] LOG:  execute <unnamed>: insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-19 10:12:12', '0', NULL, '', 'Hello', '1')\n"
	assert.Equal(t, expected, actual)
}

func TestPostgresTimestamp(t *testing.T) {
	pl := Log{}
	actual, err := pl.Timestamp("2024-04-19 10:12:16.889 CEST [89718] LOG:  execute <unnamed>: insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-19 10:12:12', '0', NULL, '', 'Hello', '1')")
	assert.Nil(t, err)
	expected, _ := time.Parse(time.DateTime, "2024-04-19 10:12:16.889")
	assert.Equal(t, expected, actual)
}

func readLine(t *testing.T, pl Log) string {
	l, err := pl.NextLine()
	if err != nil {
		t.Fatal(err)
	}
	return l
}
