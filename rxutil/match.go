package rxutil

import "regexp"

type MatchResult struct {
	rx     *regexp.Regexp
	source string
}

// Returns a MatchResult object representing the leftmost match of pattern against
// source, or nil if no matches were found.
func Match(source string, pattern string) *MatchResult {
	rx := regexp.MustCompile(pattern)

	if rx.MatchString(source) {
		return &MatchResult{
			rx:     rx,
			source: source,
		}
	}

	return nil
}

// Return the value of the numbered capture group (if given an int), or the
// named capture group (if given a string).  Returns an empty string if the
// given group name or index does not exist.
func (self *MatchResult) Group(nameOrIndex interface{}) string {
	submatches := self.rx.FindStringSubmatch(self.source)

	for i, name := range self.rx.SubexpNames() {
		switch nameOrIndex.(type) {
		case string:
			if i > 0 && name == nameOrIndex.(string) {
				return submatches[i]
			}
		case int:
			if i == nameOrIndex.(int) {
				return submatches[i]
			}
		}
	}

	return ``
}

// Returns a map of all named capture matches, keyed on capture group name.
func (self *MatchResult) NamedCaptures() map[string]string {
	captures := make(map[string]string)
	submatches := self.rx.FindStringSubmatch(self.source)

	for i, name := range self.rx.SubexpNames() {
		if i > 0 && name != `` {
			captures[name] = submatches[i]
		}
	}

	return captures
}

// Return a slice of all capture groups.
func (self *MatchResult) Captures() []string {
	return self.rx.FindStringSubmatch(self.source)
}
