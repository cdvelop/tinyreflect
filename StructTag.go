package tinyreflect

// StructTag is the tag string in a struct field (similar to reflect.StructTag)
type StructTag string

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns the empty string.
func (tag StructTag) Get(key string) string {
	value, _ := tag.Lookup(key)
	return value
}

// Lookup returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty)
// is returned. Otherwise the returned value will be the empty string.
// The ok return value reports whether the value was explicitly set in
// the tag string.
func (tag StructTag) Lookup(key string) (value string, ok bool) {
	// Simplified implementation based on Go's reflect.StructTag
	for tag != "" {
		// Skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value
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

		if key == name {
			// Unquote the value
			if len(qvalue) >= 2 && qvalue[0] == '"' && qvalue[len(qvalue)-1] == '"' {
				value = qvalue[1 : len(qvalue)-1]
				// Simple unescape for basic cases
				result := ""
				for j := 0; j < len(value); j++ {
					if value[j] == '\\' && j+1 < len(value) {
						switch value[j+1] {
						case 'n':
							result += "\n"
						case 't':
							result += "\t"
						case 'r':
							result += "\r"
						case '\\':
							result += "\\"
						case '"':
							result += "\""
						default:
							result += string(value[j])
							continue
						}
						j++ // skip the escaped character
					} else {
						result += string(value[j])
					}
				}
				return result, true
			}
			return qvalue, true
		}
	}
	return "", false
}
