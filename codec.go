package curie

import (
	"net/url"
	"unicode"
	"unicode/utf8"
)

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
