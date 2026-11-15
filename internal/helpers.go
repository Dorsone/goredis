package internal

import "strings"

func JoinValues(args []Value, sep string) string {
	var result strings.Builder

	for _, v := range args {
		switch v.GetType() {
		case STRING:
			result.WriteString(v.str + " ")
		case BULK:
			result.WriteString(v.bulk + " ")
		case ARRAY:
			result.WriteString(JoinValues(v.array, sep) + " ")
		}
	}

	return strings.Trim(result.String(), " ")
}
