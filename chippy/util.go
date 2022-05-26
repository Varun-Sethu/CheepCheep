package chippy

import "regexp"

func clearBuffer(buffer *[100]byte) {
	for i := 0; i < 100; i++ {
		buffer[i] = 0
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func matchNamedGroups(r *regexp.Regexp, value []byte) (map[string][]byte, bool) {
	match := r.FindSubmatch(value)
	if len(match) == 0 {
		return nil, false
	}

	output := make(map[string][]byte)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			output[name] = match[i]
		}
	}

	return output, true
}
