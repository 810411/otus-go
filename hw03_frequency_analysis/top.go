package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
	"unicode"
)

func Top10(s string) []string {
	const top10 = 10
	var ks []string
	cache := map[string]int{}

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && unicode.IsSymbol('-')
	}

	for _, w := range strings.FieldsFunc(s, f) {
		w = strings.ToLower(w)
		cache[w]++
	}

	for k, v := range cache {
		if v > 1 {
			ks = append(ks, k)
		}
	}

	sort.Slice(ks, func(i, j int) bool {
		return cache[ks[i]] > cache[ks[j]]
	})

	if len(ks) > top10 {
		return ks[:top10]
	}

	return ks
}
