package query

import "strings"

func getOffset(p, l int64) int64 {
	if p > 0 {
		return (p - 1) * l
	}
	return 0
}

func isPage(paramTagValue string) bool {
	return paramTagValue == "page"
}

func isLimit(paramTagValue string) bool {
	return paramTagValue == "limit"
}

func isSortBy(paramTagValue string) bool {
	if paramTagValue == "sort_by" ||
		paramTagValue == "sort-by" ||
		paramTagValue == "sortBy" ||
		paramTagValue == "sortby" {
		return true
	}
	return false
}

func validateLimit(l int64) int64 {
	if l < 1 {
		// make 10 as the default limit
		return 10
	}
	return l
}

func validatePage(p int64) int64 {
	if p < 1 {
		return 1
	}
	return p
}

// handle when all parameters is appended into single string
// e.g. []string{ "param1,param2" }
// it should convert to []string{"param1", "param2"}
func normalizeSortBy(v []string) []string {
	if len(v) == 1 {
		v[0] = strings.ReplaceAll(v[0], " ", "")
		elems := strings.Split(v[0], ",")
		v = nil
		v = append(v, elems...)
	}
	return v
}
