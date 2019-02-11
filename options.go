package gblocks

import (
	"strconv"
	"strings"
)

type Dict map[string]interface{}

func (d Dict) Contains(key string) (okay bool) {
	if _, ok := d[key]; ok {
		okay = true
	}
	return
}

// add only if the key is new
func (d Dict) Insert(key string, value interface{}) (ret interface{}) {
	if prev, ok := d[key]; !ok {
		d[key] = value
		ret = value
	} else {
		ret = prev
	}
	return
}

// make options map from tags
// based on https://golang.org/pkg/reflect/#StructTag.Get which looks up the value of a single key.
func parseTags(tag string) Dict {
	tags := make(Dict)
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		//
		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		tags[name] = parseCommas(value)
	}
	return tags
}

// returns a string or an array of strings.
func parseCommas(comma string) (ret interface{}) {
	if !strings.Contains(comma, ",") {
		ret = strings.TrimSpace(comma)
	} else {
		ret = splitCommas(comma)
	}
	return
}

// returns a comma-separated string as a slice of (trimmed) strings, or nil.
func splitCommas(comma string) (ret interface{}) {
	if ar := strings.Split(comma, ","); len(ar) > 0 {
		for i, s := range ar {
			ar[i] = strings.TrimSpace(s)
		}
		ret = ar
	}
	return
}
