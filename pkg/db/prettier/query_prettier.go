package prettier

import (
	"fmt"
	"strings"
)

// Placeholders block
const (
	PlaceholderDollar = "$"
)

// Pretty sql query
func Pretty(query, placeholder string, args ...any) string {
	for i, param := range args {
		var val string
		switch v := param.(type) {
		case string:
			val = fmt.Sprintf("%q", v)
		case []byte:
			val = fmt.Sprintf("%q", string(v))
		default:
			val = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, fmt.Sprintf("%s%d", placeholder, i+1), val, -1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}
