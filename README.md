# dinheiro

A Go library for handling Brazilian monetary values.

## Installation

```bash
go get github.com/mvfavila/dinheiro
```

## Overview

Values are represented as `int64` centavos internally. Both functions accept equivalent inputs:

| int64 (centavos) | string (raw centavos) | string (formatted) |
|------------------|-----------------------|--------------------|
| `2`              | `"2"`                 | `"0,02"`           |
| `199`            | `"199"`               | `"1,99"`           |
| `100150`         | `"100150"`            | `"1001,50"`       |
| `100137`         | `"100137"`            | `"1.001,37"`       |

Negative values and unsupported types return an error.

## Functions

### `ToText(value any) (string, error)`

Formats a monetary value as a Brazilian currency string (dots as thousand separators, comma as decimal separator).

```go
dinheiro.ToText(int64(2))          // "0,02"
dinheiro.ToText("30")              // "0,30"
dinheiro.ToText("1,99")            // "1,99"
dinheiro.ToText(int64(100137))     // "1.001,37"
dinheiro.ToText("7700022280")      // "77.000.222,80"
dinheiro.ToText("77000222,80")     // "77.000.222,80"
dinheiro.ToText("77.000.222,80")   // "77.000.222,80"
dinheiro.ToText("77.000.222")      // error
```

### `ToTextDescription(value any) (string, error)`

Formats a monetary value as a written-out Brazilian Portuguese description.

```go
dinheiro.ToTextDescription(int64(2))         // "dois centavos"
dinheiro.ToTextDescription("30")             // "trinta centavos"
dinheiro.ToTextDescription("1,99")           // "um real e noventa e nove centavos"
dinheiro.ToTextDescription(int64(100137))    // "um mil e um reais e trinta e sete centavos"
dinheiro.ToTextDescription("7700022280")     // "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"
dinheiro.ToTextDescription("77000222,80")    // "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"
dinheiro.ToTextDescription("77.000.222,80")  // "setenta e sete milhões duzentos e vinte e dois reais e oitenta centavos"
dinheiro.ToTextDescription("77.000.222")     // error
```

## Supported range

Any non-negative `int64` value (up to ~92 quadrilhões reais).

## License

See [LICENSE](LICENSE).
