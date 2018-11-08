package httputil

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
)

// Parses the named query string from a request as an integer.
func QInt(req *http.Request, key string, fallbacks ...int64) int64 {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToInteger(v); err == nil {
			return i
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return 0
	}
}

// Parses the named query string from a request as a float.
func QFloat(req *http.Request, key string, fallbacks ...float64) float64 {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToFloat(v); err == nil {
			return i
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return 0
	}
}

// Parses the named query string from a request as a date/time value.
func QTime(req *http.Request, key string) time.Time {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToTime(v); err == nil {
			return i
		}
	}

	return time.Time{}
}

// Parses the named query string from a request as a boolean value.
func QBool(req *http.Request, key string, fallbacks ...bool) bool {
	if v := Q(req, key); v == `` && len(fallbacks) > 0 {
		return fallbacks[0]
	} else if typeutil.Bool(v) {
		return true
	}

	return false
}

// Parses the named query string from a request as a string.
func Q(req *http.Request, key string, fallbacks ...string) string {
	if v := req.URL.Query().Get(key); v != `` {
		if vS, err := url.QueryUnescape(v); err == nil {
			return vS
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return ``
	}
}

// Parses the named query string from a request as a delimiter-separated string slice.
func QStrings(req *http.Request, key string, delimiter string, fallbacks ...string) []string {
	if strs := sliceutil.CompactString(strings.Split(Q(req, key), delimiter)); len(strs) > 0 {
		return strs
	} else if len(fallbacks) > 0 {
		return sliceutil.Stringify(sliceutil.Flatten(fallbacks))
	} else {
		return make([]string, 0)
	}
}

// Sets a query string to the given value in the given url.URL
func SetQ(u *url.URL, key string, value interface{}) {
	qs := u.Query()
	qs.Set(key, stringutil.MustString(value))
	u.RawQuery = qs.Encode()
}

// Appends a query string from then given url.URL
func AddQ(u *url.URL, key string, value interface{}) {
	qs := u.Query()
	qs.Add(key, stringutil.MustString(value))
	u.RawQuery = qs.Encode()
}

// Deletes a query string from then given url.URL
func DelQ(u *url.URL, key string) {
	qs := u.Query()
	qs.Del(key)
	u.RawQuery = qs.Encode()
}
