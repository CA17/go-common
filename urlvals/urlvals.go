package urlvals

import (
	"net/url"
	"strconv"
)

func GetInt64Value(vals url.Values, name string, defval int64) int64 {
	v := vals.Get(name)
	if v == "" {
		return defval
	}
	vv, _ := strconv.ParseInt(v, 10, 64)
	return vv
}

func GetIntValue(vals url.Values, name string, defval int) int {
	v := vals.Get(name)
	if v == "" {
		return defval
	}
	vv, _ := strconv.Atoi(v)
	return vv
}

func GetStrValue(vals url.Values, name string, defval string) string {
	src := vals.Get(name)
	if src == "" {
		return defval
	}
	return src
}
