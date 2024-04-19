package matcher

// PostgresTokenizer tokenizes PostgreSQL log entries.
// PostgreSQL configuration settings:
//
//	log_destination = 'stderr,csvlog'
//
// PostgreSQL splits single sql statements into two parts:
//
//	1: select * from job where job0_.publish_at<$1 and job0_.publish_trials<$2
//	2: parameters: $1 = '2024-04-19 10:07:38.543981', $2 = '1'
type PostgresTokenizer struct {
}

func (m PostgresTokenizer) Tokenize(s string, patterns []string) []string {
	return Tokenize(normalize(cutPrefix(s, patterns), patterns))
}
