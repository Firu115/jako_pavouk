package utils

import (
	"regexp"
	"strings"
)

var uvozovka string = `"`
var pomlcka string = `-`
var tabulka = map[rune]string{
	'„': uvozovka,
	'”': uvozovka,
	'“': uvozovka,
	'‟': uvozovka,
	'❞': uvozovka,
	'❝': uvozovka,
	'🙷': uvozovka,
	'🙶': uvozovka,
	'＂': uvozovka,
	'‚': uvozovka,
	'’': uvozovka,
	'‘': uvozovka,
	'‛': uvozovka,
	'❛': uvozovka,
	'❜': uvozovka,
	'⹂': uvozovka,
	'…': "...",
	'‑': pomlcka,
	'‐': pomlcka,
	'־': pomlcka,
	'‒': pomlcka,
	'–': pomlcka,
	'—': pomlcka,
	'―': pomlcka,
	'⸺': pomlcka,
	'⸻': pomlcka,
	'﹘': pomlcka,
	'﹣': pomlcka,
	'－': pomlcka,
	'֊': pomlcka,
	'᠆': pomlcka,
	' ': " ",
	'ä': "a", // nemecke
	'ö': "o",
	'ü': "u",
	'ć': "c", // polske
	'ą': "a",
	'ę': "e",
	'ł': "l",
	'ń': "n",
	'ó': "o",
	'ź': "z",
	'ż': "z",
	'ĺ': "l", //slovenske
	'ľ': "l",
	'ô': "o",
	'ŕ': "r",
	'è': "e", // ani nevim
	'ã': "a",
	'ë': "e",
	'æ': "ae",
}
var tabulka2 map[string]string = map[string]string{
	"…": "...",
}
var mezery = regexp.MustCompile("( {2,})")

func UpravaTextu(txt string) string {
	var v strings.Builder
	for _, ch := range txt {
		if tabulka[ch] != "" {
			v.WriteString(tabulka[ch])
		} else {
			v.WriteRune(ch)
		}
	}

	for o, n := range tabulka2 {
		txt = strings.ReplaceAll(txt, o, n)
	}

	vysledek := v.String()
	vysledek = mezery.ReplaceAllString(vysledek, " ") // odstrani vice mezer

	return strings.TrimSpace(vysledek)
}
