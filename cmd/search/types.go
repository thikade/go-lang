package search

// need to define search and result obj becausewe use the same html template
// and template code will throw errors if it cannot find "TotalResults" in obj
type SearchObj struct {
	Days   string
	Token  string
	Errors map[string]string
	// search response
	TotalResults int
	Results      []string
}
