//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

package curie

import (
	"testing"

	"github.com/fogfish/it"
)

func TestDecode(t *testing.T) {
	for uri, iri := range map[string]string{
		"wiki":                              "wiki",
		"%E1%BF%AC%CF%8C%CE%B4%CE%BF%CF%82": "Ῥόδος",
		"%e1%bf%ac%cf%8c%ce%b4%ce%bf%cf%82": "Ῥόδος",
		"%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BA%D1%82%D0%BD%D1%8B%D0%B5%20%D1%83%D0%BD%D0%B8%D1%84%D0%B8%D1%86%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%BD%D1%8B%D0%B5%20%D0%B8%D0%B4%D0%B5%D0%BD%D1%82%D0%B8%D1%84%D0%B8%D0%BA%D0%B0%D1%82%D0%BE%D1%80%D1%8B%20%D1%80%D0%B5%D1%81%D1%83%D1%80%D1%81%D0%BE%D0%B2": "Компактные унифицированные идентификаторы ресурсов",
		"%00%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F": "%00%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F",
		"-._~":               "-._~",
		"%2D%2E%5F%7E":       "-._~",
		":/?#[]@!$&'()*+,;=": ":/?#[]@!$&'()*+,;=",
		"%3A%2F%3F%23%5B%5D%40%21%24%26%27%28%29%2A%2B%2C%3B%3D": "%3A%2F%3F%23%5B%5D%40%21%24%26%27%28%29%2A%2B%2C%3B%3D",
		"%%%":      "%%%",
		"%Ww%wW%%": "%Ww%wW%%",
		"%s":       "%s",
	} {
		it.Ok(t).If(Decode(uri)).Equal(iri)
	}
}

func TestUnHex(t *testing.T) {
	for hex, val := range map[byte]byte{
		0x35: 0x05,
		0x63: 0x0c,
		0x43: 0x0c,
		0x00: 0x00,
		0x4F: 0x00,
	} {
		it.Ok(t).If(unhex(hex)).Equal(val)
	}
}
