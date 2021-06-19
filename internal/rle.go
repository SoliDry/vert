package internal

import (
	"strconv"
)

func RunLengthEncode(s string) string {
	e := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		j := i + 1
		for ; j <= len(s); j++ {
			if j < len(s) && s[j] == c {
				continue
			}
			if j-i > 1 {
				e = strconv.AppendInt(e, int64(j-i), 10)
			}
			e = append(e, c)
			break
		}
		i = j - 1
	}

	return string(e)
}

func RunLengthDecode(s string) string {
	d := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		n := 0
		for ; i < len(s) && (s[i] >= '0' && s[i] <= '9'); i++ {
			n = 10*n + int(s[i]-'0')
		}
		if i < len(s) {
			c := s[i]
			for ; n > 1; n-- {
				d = append(d, c)
			}
			d = append(d, c)
		}
	}

	return string(d)
}
