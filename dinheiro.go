// Package dinheiro provides utilities for handling Brazilian monetary values.
// Values are represented as int64 centavos (smallest unit).
package dinheiro

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// parseInput converts the input value to a non-negative int64 representing centavos.
// Accepted types: int64 (raw centavos) or string (raw centavos like "199", or
// Brazilian-formatted like "1,99" or "1.001,37").
func parseInput(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		if v < 0 {
			return 0, errors.New("dinheiro: negative values are not accepted")
		}
		return v, nil
	case string:
		n, err := parseString(v)
		if err != nil {
			return 0, err
		}
		if n < 0 {
			return 0, errors.New("dinheiro: negative values are not accepted")
		}
		return n, nil
	default:
		return 0, fmt.Errorf("dinheiro: unsupported type %T; only int64 and string are accepted", value)
	}
}

// parseString parses a string into centavos.
// Accepts raw integers ("2", "100137") or Brazilian-formatted ("0,02", "1001,50", "1.001,37").
func parseString(s string) (int64, error) {
	if strings.ContainsRune(s, ',') {
		return parseFormattedString(s)
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("dinheiro: invalid value %q: %w", s, err)
	}
	return n, nil
}

// parseFormattedString parses a Brazilian-formatted monetary string into centavos.
// The format is dots as thousand separators and a comma as the decimal separator.
// Example: "1.001,37" → 100137.
func parseFormattedString(s string) (int64, error) {
	invalid := func() (int64, error) {
		return 0, fmt.Errorf("dinheiro: invalid monetary format %q; expected format like \"1.234,56\"", s)
	}

	// Remove thousand separators (dots).
	cleaned := strings.ReplaceAll(s, ".", "")

	parts := strings.SplitN(cleaned, ",", 2)
	if len(parts) != 2 || len(parts[1]) != 2 {
		return invalid()
	}

	reais, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || reais < 0 {
		return invalid()
	}

	centavos, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || centavos < 0 {
		return invalid()
	}

	// Guard against overflow: reais*100 + centavos must fit in int64.
	if reais > (math.MaxInt64-centavos)/100 {
		return 0, fmt.Errorf("dinheiro: value too large")
	}

	return reais*100 + centavos, nil
}

// formatWithDots formats a non-negative int64 with dot thousand separators.
// Example: 1001 → "1.001", 77000222 → "77.000.222".
func formatWithDots(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	buf := make([]byte, 0, len(s)+(len(s)-1)/3)
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			buf = append(buf, '.')
		}
		buf = append(buf, s[i])
	}
	return string(buf)
}

// ToText formats a monetary value as a Brazilian currency string.
//
// The input may be an int64 (raw centavos) or a string in either raw-centavos
// format ("199") or Brazilian-formatted currency ("1,99" or "1.001,37").
// All three representations are equivalent.
// Negative values are rejected.
//
// Examples:
//
//	ToText(int64(2))        → "0,02"
//	ToText("30")            → "0,30"
//	ToText("1,99")          → "1,99"
//	ToText(int64(100137))   → "1.001,37"
//	ToText("77.000.222,80") → "77.000.222,80"
func ToText(value any) (string, error) {
	n, err := parseInput(value)
	if err != nil {
		return "", err
	}
	reais := n / 100
	centavos := n % 100
	return fmt.Sprintf("%s,%02d", formatWithDots(reais), centavos), nil
}

// ToTextDescription formats a monetary value as a written-out Brazilian
// Portuguese description.
//
// The input may be an int64 (raw centavos) or a string in either raw-centavos
// format ("199") or Brazilian-formatted currency ("1,99" or "1.001,37").
// All three representations are equivalent.
// Negative values are rejected.
//
// Examples:
//
//	ToTextDescription(int64(2))          → "dois centavos"
//	ToTextDescription("30")              → "trinta centavos"
//	ToTextDescription("1,99")            → "um real e noventa e nove centavos"
//	ToTextDescription(int64(100137))     → "um mil e um reais e trinta e sete centavos"
//	ToTextDescription("77.000.222,80")   → "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"
func ToTextDescription(value any) (string, error) {
	n, err := parseInput(value)
	if err != nil {
		return "", err
	}

	if n == 0 {
		return "zero centavos", nil
	}

	reais := n / 100
	centavos := n % 100

	var parts []string

	if reais > 0 {
		w := numberToWords(reais)
		if reais == 1 {
			parts = append(parts, w+" real")
		} else {
			parts = append(parts, w+" reais")
		}
	}

	if centavos > 0 {
		w := numberToWords(centavos)
		if centavos == 1 {
			parts = append(parts, w+" centavo")
		} else {
			parts = append(parts, w+" centavos")
		}
	}

	return strings.Join(parts, " e "), nil
}
