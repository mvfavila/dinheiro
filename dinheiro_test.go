package dinheiro_test

import (
	"math"
	"testing"

	"github.com/mvfavila/dinheiro"
)

// --- ToText ---

func TestToText_Int64(t *testing.T) {
	cases := []struct {
		input    int64
		expected string
	}{
		{0, "0,00"},
		{1, "0,01"},
		{2, "0,02"},
		{30, "0,30"},
		{99, "0,99"},
		{100, "1,00"},
		{101, "1,01"},
		{199, "1,99"},
		{1000, "10,00"},
		{100137, "1.001,37"},
		{7700022280, "77.000.222,80"},
		{1000000000, "10.000.000,00"},
	}

	for _, tc := range cases {
		got, err := dinheiro.ToText(tc.input)
		if err != nil {
			t.Errorf("ToText(%d): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ToText(%d) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestToText_StringRaw(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"0", "0,00"},
		{"2", "0,02"},
		{"30", "0,30"},
		{"199", "1,99"},
		{"100137", "1.001,37"},
		{"7700022280", "77.000.222,80"},
	}

	for _, tc := range cases {
		got, err := dinheiro.ToText(tc.input)
		if err != nil {
			t.Errorf("ToText(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ToText(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestToText_StringFormatted(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"0,00", "0,00"},
		{"0,02", "0,02"},
		{"0,30", "0,30"},
		{"1001,50", "1.001,50"},
		{"1,99", "1,99"},
		{"1.001,37", "1.001,37"},
		{"77.000.222,80", "77.000.222,80"},
	}

	for _, tc := range cases {
		got, err := dinheiro.ToText(tc.input)
		if err != nil {
			t.Errorf("ToText(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ToText(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

// Equivalent representations must produce the same output.
func TestToText_Equivalence(t *testing.T) {
	cases := []struct {
		int64Input   int64
		rawStr       string
		formattedStr string
		expected     string
	}{
		{2, "2", "0,02", "0,02"},
		{30, "30", "0,30", "0,30"},
		{199, "199", "1,99", "1,99"},
		{100137, "100137", "1.001,37", "1.001,37"},
		{7700022280, "7700022280", "77.000.222,80", "77.000.222,80"},
	}

	for _, tc := range cases {
		r1, err := dinheiro.ToText(tc.int64Input)
		if err != nil {
			t.Errorf("ToText(int64(%d)): unexpected error: %v", tc.int64Input, err)
			continue
		}
		r2, err := dinheiro.ToText(tc.rawStr)
		if err != nil {
			t.Errorf("ToText(%q): unexpected error: %v", tc.rawStr, err)
			continue
		}
		r3, err := dinheiro.ToText(tc.formattedStr)
		if err != nil {
			t.Errorf("ToText(%q): unexpected error: %v", tc.formattedStr, err)
			continue
		}
		if r1 != tc.expected || r2 != tc.expected || r3 != tc.expected {
			t.Errorf("equivalence failed for value %d: int64=%q raw=%q formatted=%q want %q",
				tc.int64Input, r1, r2, r3, tc.expected)
		}
	}
}

func TestToText_Errors(t *testing.T) {
	_, err := dinheiro.ToText(int64(-1))
	if err == nil {
		t.Error("ToText(int64(-1)): expected error for negative value")
	}

	_, err = dinheiro.ToText("-1")
	if err == nil {
		t.Error(`ToText("-1"): expected error for negative string`)
	}

	_, err = dinheiro.ToText("abc")
	if err == nil {
		t.Error(`ToText("abc"): expected error for non-numeric string`)
	}

	_, err = dinheiro.ToText(float64(1.99))
	if err == nil {
		t.Error("ToText(float64(1.99)): expected error for unsupported type")
	}

	_, err = dinheiro.ToText(42) // int, not int64
	if err == nil {
		t.Error("ToText(42): expected error for int (not int64)")
	}
}

// --- ToMoneyDescription ---

func TestToMoneyDescription_Int64(t *testing.T) {
	cases := []struct {
		input    int64
		expected string
	}{
		{0, "zero centavos"},
		{1, "um centavo"},
		{2, "dois centavos"},
		{30, "trinta centavos"},
		{99, "noventa e nove centavos"},
		{100, "um real"},
		{101, "um real e um centavo"},
		{199, "um real e noventa e nove centavos"},
		{200, "dois reais"},
		{1000, "dez reais"},
		{10000, "cem reais"},
		{100137, "um mil e um reais e trinta e sete centavos"},
		{7700022280, "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"},
		{500, "cinco reais"},
		{math.MaxInt64, "noventa e dois quatrilhões duzentos e trinta e três trilhões setecentos e vinte bilhões trezentos e sessenta e oito milhões quinhentos e quarenta e sete mil setecentos e cinquenta e oito reais e sete centavos"},
	}

	for _, tc := range cases {
		got, err := dinheiro.ToMoneyDescription(tc.input)
		if err != nil {
			t.Errorf("ToMoneyDescription(%d): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ToMoneyDescription(%d) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestToMoneyDescription_StringRaw(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"2", "dois centavos"},
		{"30", "trinta centavos"},
		{"1001,50", "um mil e um reais e cinquenta centavos"},
		{"199", "um real e noventa e nove centavos"},
		{"100137", "um mil e um reais e trinta e sete centavos"},
		{"7700022280", "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"},
	}

	for _, tc := range cases {
		got, err := dinheiro.ToMoneyDescription(tc.input)
		if err != nil {
			t.Errorf("ToMoneyDescription(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ToMoneyDescription(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestToMoneyDescription_StringFormatted(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"0,02", "dois centavos"},
		{"0,30", "trinta centavos"},
		{"1,99", "um real e noventa e nove centavos"},
		{"1.001,37", "um mil e um reais e trinta e sete centavos"},
		{"77.000.222,80", "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"},
	}

	for _, tc := range cases {
		got, err := dinheiro.ToMoneyDescription(tc.input)
		if err != nil {
			t.Errorf("ToMoneyDescription(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ToMoneyDescription(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestToMoneyDescription_Equivalence(t *testing.T) {
	cases := []struct {
		int64Input   int64
		rawStr       string
		formattedStr string
		expected     string
	}{
		{2, "2", "0,02", "dois centavos"},
		{30, "30", "0,30", "trinta centavos"},
		{199, "199", "1,99", "um real e noventa e nove centavos"},
		{100137, "100137", "1.001,37", "um mil e um reais e trinta e sete centavos"},
		{7700022280, "7700022280", "77.000.222,80", "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"},
	}

	for _, tc := range cases {
		r1, err := dinheiro.ToMoneyDescription(tc.int64Input)
		if err != nil {
			t.Errorf("ToMoneyDescription(int64(%d)): unexpected error: %v", tc.int64Input, err)
			continue
		}
		r2, err := dinheiro.ToMoneyDescription(tc.rawStr)
		if err != nil {
			t.Errorf("ToMoneyDescription(%q): unexpected error: %v", tc.rawStr, err)
			continue
		}
		r3, err := dinheiro.ToMoneyDescription(tc.formattedStr)
		if err != nil {
			t.Errorf("ToMoneyDescription(%q): unexpected error: %v", tc.formattedStr, err)
			continue
		}
		if r1 != tc.expected || r2 != tc.expected || r3 != tc.expected {
			t.Errorf("equivalence failed for value %d: int64=%q raw=%q formatted=%q want %q",
				tc.int64Input, r1, r2, r3, tc.expected)
		}
	}
}

func TestToMoneyDescription_Errors(t *testing.T) {
	_, err := dinheiro.ToMoneyDescription(int64(-1))
	if err == nil {
		t.Error("ToMoneyDescription(int64(-1)): expected error for negative value")
	}

	_, err = dinheiro.ToMoneyDescription("-100")
	if err == nil {
		t.Error(`ToMoneyDescription("-100"): expected error for negative string`)
	}

	_, err = dinheiro.ToMoneyDescription("xyz")
	if err == nil {
		t.Error(`ToMoneyDescription("xyz"): expected error for non-numeric string`)
	}

	_, err = dinheiro.ToMoneyDescription(3.14)
	if err == nil {
		t.Error("ToMoneyDescription(3.14): expected error for unsupported type float64")
	}
}

// --- "e" connector rule tests ---

func TestToMoneyDescription_EConnector(t *testing.T) {
	cases := []struct {
		input    int64
		expected string
	}{
		// remainder <= 100 -> "e"
		{100100, "um mil e um reais"},
		{110000, "um mil e cem reais"},
		// remainder > 100 -> no "e"
		{110100, "um mil cento e um reais"},
		{120000, "um mil duzentos reais"},
	}

	for _, tc := range cases {
		got, err := dinheiro.ToMoneyDescription(tc.input)
		if err != nil {
			t.Errorf("ToMoneyDescription(%d): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ToMoneyDescription(%d) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}
