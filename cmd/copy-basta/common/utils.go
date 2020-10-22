package common

import "strings"

func TrimRootDir(s string) string {
	ss := strings.Split(s, "/")
	if len(ss) == 1 {
		return ss[0]
	}
	return strings.Join(ss[1:], "/")
}
