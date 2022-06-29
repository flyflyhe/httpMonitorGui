package stringHelper

import "strings"

func SplitStringByLength(str, sep string, l int) string {
	strLen := len(str)
	if strLen <= l {
		return str
	}
	r := []rune(str)

	sb := strings.Builder{}
	for i := 0; i < len(r); i++ {
		if i > 1 && i%l == 0 {
			sb.WriteString(sep)
		}
		sb.WriteRune(r[i])
	}

	return sb.String()
}
