package tin

import (
	"strconv"
	"strings"

	"github.com/ionous/gblocks/block"
)

// make dict map from tags
// based on https://golang.org/pkg/reflect/#StructTag.Get which looks up the value of a single key.
func parseTags(tag string) (ret block.Dict) {
	if len(tag) > 0 {
		tags := make(block.Dict)
		visitTags(tag, func(k, v string) {
			tags[k] = parseCommas(v)
		})
		ret = tags
	}
	return
}

func visitTags(tag string, cb func(k, v string)) {
	for len(tag) > 0 {
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
		cb(name, value)
	}
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
