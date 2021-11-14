package git

import (
	"regexp"
	"strings"
)

func shellQuote(s string) string {
	switch {
	case s == "":
		return "''"
	case safeRE.MatchString(s):
		return s
	case strings.IndexByte(s, '\'') == -1:
		return "'" + s + "'"
	default:
		var buf strings.Builder
		buf.WriteByte('"')
		for _, b := range []byte(s) {
			switch {
			case charBits[b]&bitEscape != 0:
				buf.WriteByte('\\')
				buf.WriteByte('c')
			case b < 0x20 || b > 0x7F:
				buf.WriteByte('\\')
				buf.WriteByte('x')
				buf.WriteByte(hexChars[(b>>4)&0xF])
				buf.WriteByte(hexChars[b&0xF])
			default:
				buf.WriteByte(b)
			}
		}
		buf.WriteByte('"')
		return buf.String()
	}
}

const (
	safeChars      = "-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/.@_+=%^:"
	escapableChars = "$\"\\`"
)

var (
	// TODO: find all safe characters
	safeRE = regexp.MustCompile("^[" + safeChars + "]+$")

	charBits [256]charBit

	hexChars = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
)

type charBit int

const (
	bitSafe charBit = 1 << iota
	bitEscape
)

func init() {
	for _, c := range safeChars {
		charBits[c] |= bitSafe
	}
	for _, c := range []byte{'$', '"', '\\', '`', '!'} {
		charBits[c] |= bitEscape
	}
}
