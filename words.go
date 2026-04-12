package dinheiro

var unitWords = []string{
	"", "um", "dois", "três", "quatro", "cinco", "seis", "sete", "oito", "nove",
	"dez", "onze", "doze", "treze", "quatorze", "quinze", "dezesseis", "dezessete", "dezoito", "dezenove",
}

var tenWords = []string{
	"", "", "vinte", "trinta", "quarenta", "cinquenta", "sessenta", "setenta", "oitenta", "noventa",
}

var hundredWords = []string{
	"", "cento", "duzentos", "trezentos", "quatrocentos", "quinhentos",
	"seiscentos", "setecentos", "oitocentos", "novecentos",
}

// belowHundred returns Brazilian Portuguese words for 1–99.
func belowHundred(n int64) string {
	if n < 20 {
		return unitWords[n]
	}
	t := tenWords[n/10]
	if u := n % 10; u != 0 {
		return t + " e " + unitWords[u]
	}
	return t
}

// belowThousand returns Brazilian Portuguese words for 1–999.
func belowThousand(n int64) string {
	if n < 100 {
		return belowHundred(n)
	}
	h := n / 100
	rem := n % 100
	if rem == 0 {
		if h == 1 {
			return "cem"
		}
		return hundredWords[h]
	}
	return hundredWords[h] + " e " + belowHundred(rem)
}

// numberToWords converts a non-negative int64 to Brazilian Portuguese words.
// The "e" connector between groups is only used when the total value of all
// remaining groups is ≤ 100 (e.g. "um mil e um", but "um mil cento e um").
func numberToWords(n int64) string {
	if n == 0 {
		return "zero"
	}

	type chunk struct {
		text  string
		value int64
	}

	groups := []struct {
		divisor  int64
		singular string
		plural   string
	}{
		{1_000_000_000_000_000, "um quatrilhão", "quatrilhões"},
		{1_000_000_000_000, "um trilhão", "trilhões"},
		{1_000_000_000, "um bilhão", "bilhões"},
		{1_000_000, "um milhão", "milhões"},
		{1_000, "", "mil"},
	}

	var chunks []chunk
	rem := n

	for _, g := range groups {
		if rem < g.divisor {
			continue
		}
		count := rem / g.divisor
		rem -= count * g.divisor

		var text string
		switch {
		case g.divisor == 1_000:
			// Always "X mil" (e.g. "um mil", "dois mil", "cem mil")
			text = belowThousand(count) + " mil"
		case count == 1:
			text = g.singular
		default:
			text = belowThousand(count) + " " + g.plural
		}

		chunks = append(chunks, chunk{text, count * g.divisor})
	}

	if rem > 0 {
		chunks = append(chunks, chunk{belowThousand(rem), rem})
	}

	result := chunks[0].text
	for i := 1; i < len(chunks); i++ {
		// Sum values of remaining chunks to decide the connector.
		var tail int64
		for j := i; j < len(chunks); j++ {
			tail += chunks[j].value
		}
		if tail <= 100 {
			result += " e " + chunks[i].text
		} else {
			result += " " + chunks[i].text
		}
	}

	return result
}
