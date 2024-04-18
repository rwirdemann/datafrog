package matcher

// A Tokenizer splits a raw string obtained from the log file of the monitored
// API into a set of single tokens. The primary split character should be a
// space. All leading and trailing spaces should be removed from each token. The
// first token should always be one of the matching patterns provided to
// Tokenize. Thus, all leading timestamps, etc. should be cut from the raw
// string before the tokenizing begins. Example: The raw string
// "2024-04-08T09:39:15.070009Z	 2549 Query	insert into job (description,
// id) values ('Developer', 5)" should become the following token set:
// ["insert", "into", "job", "(description,", "id)", "values", "('Developer',",
// "5)"]
type Tokenizer interface {
	Tokenize(s string, patterns []string) []string
}
