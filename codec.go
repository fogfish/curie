package curie

import (
	"net/url"
	"unicode"
	"unicode/utf8"
)

const upperHex = "0123456789ABCDEF"

// Decode converts URIs to IRIs as defined by RFC 3987
// https://www.rfc-editor.org/rfc/rfc3987#section-3.2
func Decode(uri string) string {
	return string(decode([]byte(uri)))
}

func decode(uri []byte) []byte {
	var iri []byte

	decodeLastRune := func() {
		r, size := utf8.DecodeLastRune(iri)
		if !unicode.IsGraphic(r) && r != utf8.RuneError {
			esc := url.PathEscape(string(r))
			iri = append(iri[:len(iri)-size], []byte(esc)...)
		}
	}

	for i := 0; i < len(uri); {
		switch {
		case uri[i] == '%' && i+2 <= len(uri) && ishex(uri[i+1]) && ishex(uri[i+2]):
			b := unhex(uri[i+1])<<4 | unhex(uri[i+2])
			if checkReserved(b) {
				iri = append(iri, uri[i:i+3]...)
			} else {
				iri = append(iri, b)
				decodeLastRune()
			}
			i += 3
		default:
			iri = append(iri, uri[i])
			decodeLastRune()
			i++
		}
	}

	return iri
}

func checkReserved(b byte) bool {
	return b == ':' ||
		b == '/' || b == '?' || b == '#' ||
		b == '[' || b == ']' || b == '@' ||
		b == '!' || b == '$' || b == '&' ||
		b == '\'' || b == '(' || b == ')' ||
		b == '*' || b == '+' || b == ',' ||
		b == ';' || b == '='
}

// Reserved characters are:
//
//	%21 %23 %24 %26 %27 %28 %29 %2A %2B %2C %2F
//	%3A %3B %3D %3F
//	%40
//	%5B %5D
// func isReserved(hi, lo byte) bool {
// 	switch {
// 	case hi == '4' && lo == '0':
// 		return true
// 	case hi == '5' && (lo == 'B' || lo == 'b' || lo == 'D' || lo == 'd'):
// 		return true
// 	case hi == '3' && (lo == 'A' || lo == 'a' || lo == 'B' || lo == 'b' || lo == 'D' || lo == 'd' || lo == 'F' || lo == 'f'):
// 		return true
// 	case hi == '2' && (lo == '1' || lo == '3' || lo == '4' || lo == '6' || lo == '7' || lo == '8' || lo == '9'):
// 		return true
// 	case hi == '2' && (lo == 'A' || lo == 'a' || lo == 'B' || lo == 'b' || lo == 'C' || lo == 'c' || lo == 'F' || lo == 'f'):
// 		return true
// 	default:
// 		return false
// 	}
// }

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}
