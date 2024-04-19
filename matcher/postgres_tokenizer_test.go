package matcher

import (
	"fmt"
	"testing"
)

func TestTokenizePostgresLog(t *testing.T) {
	pt := PostgresTokenizer{}
	s := "select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ where (job0_.published_timestamp is null) and job0_.publish_at<$1 and job0_.publish_trials<$2"
	tokens := pt.Tokenize(s, []string{"select"})
	fmt.Printf("%v\n", tokens)
}
