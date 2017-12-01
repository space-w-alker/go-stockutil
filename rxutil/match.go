// Utilities to make working with regular expressions easier.
package rxutil

import "regexp"

type MatchResult struct {
	rx                  *regexp.Regexp
	source              string
	leftmost_submatches []string
	all_submatches      [][]string
}

// Returns a MatchResult object representing the leftmost match of pattern against
// source, or nil if no matches were found.
func Match(source string, pattern string) *MatchResult {
	rx := regexp.MustCompile(pattern)

	if rx.MatchString(source) {
		return &MatchResult{
			rx:                  rx,
			source:              source,
			leftmost_submatches: rx.FindStringSubmatch(source),
			all_submatches:      rx.FindAllStringSubmatch(source, -1),
		}
	}

	return nil
}

// Return the value of the numbered capture group (if given an int), or the
// named capture group (if given a string).  Returns an empty string if the
// given group name or index does not exist.
func (self *MatchResult) Group(nameOrIndex interface{}) string {
	for i, name := range self.rx.SubexpNames() {
		switch nameOrIndex.(type) {
		case string:
			if i > 0 && name == nameOrIndex.(string) {
				return self.leftmost_submatches[i]
			}
		case int:
			if i == nameOrIndex.(int) {
				return self.leftmost_submatches[i]
			}
		}
	}

	return ``
}

// Returns a map of all named capture matches, keyed on capture group name.
func (self *MatchResult) NamedCaptures() map[string]string {
	captures := make(map[string]string)

	for i, name := range self.rx.SubexpNames() {
		if i > 0 && name != `` {
			captures[name] = self.leftmost_submatches[i]
		}
	}

	return captures
}

// Return a slice of all capture groups.
func (self *MatchResult) Captures() []string {
	return self.leftmost_submatches
}

// Returns all captures from all matches appended together.  The full match string
// from match is omitted, so only the actual values appearing within capture groups
// are returned.
func (self *MatchResult) AllCaptures() []string {
	all := make([]string, 0)

	for _, subcaptures := range self.all_submatches {
		if len(subcaptures) > 1 {
			all = append(all, subcaptures[1:]...)
		}
	}

	return all
}
