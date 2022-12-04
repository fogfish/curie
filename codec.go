package curie

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const upperHex = "0123456789ABCDEF"

// Decode converts URIs to IRIs as defined by RFC 3987
// https://www.rfc-editor.org/rfc/rfc3987#section-3.2
func Decode(uri string) string {
	var iri strings.Builder

	for i := 0; i < len(uri); {
		switch {
		case uri[i] == '%' && i+2 <= len(uri) && ishex(uri[i+1]) && ishex(uri[i+2]):
			hi, lo := uri[i+1], uri[i+2]
			if isReserved(hi, lo) {
				iri.WriteString(uri[i : i+3])
			} else {
				iri.WriteByte(unhex(hi)<<4 | unhex(lo))
			}
			i += 3
		default:
			iri.WriteByte(uri[i])
			i++
		}
	}

	raw := iri.String()

	var str strings.Builder
	for index, r := range raw {
		switch {
		case r == utf8.RuneError:
			str.WriteString("%FF%FD")
		case unicode.IsGraphic(r):
			str.WriteRune(r)
		default:
			b := raw[index]
			str.WriteByte('%')
			str.WriteByte(upperHex[b>>4])
			str.WriteByte(upperHex[b&0xF])
		}
	}
	return str.String()
}

// Reserved characters are:
//
//	%21 %23 %24 %26 %27 %28 %29 %2A %2B %2C %2F
//	%3A %3B %3D %3F
//	%40
//	%5B %5D
func isReserved(hi, lo byte) bool {
	switch {
	case hi == '4' && lo == '0':
		return true
	case hi == '5' && (lo == 'B' || lo == 'b' || lo == 'D' || lo == 'd'):
		return true
	case hi == '3' && (lo == 'A' || lo == 'a' || lo == 'B' || lo == 'b' || lo == 'D' || lo == 'd' || lo == 'F' || lo == 'f'):
		return true
	case hi == '2' && (lo == '1' || lo == '3' || lo == '4' || lo == '6' || lo == '7' || lo == '8' || lo == '9'):
		return true
	case hi == '2' && (lo == 'A' || lo == 'a' || lo == 'B' || lo == 'b' || lo == 'C' || lo == 'c' || lo == 'F' || lo == 'f'):
		return true
	default:
		return false
	}
}

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
