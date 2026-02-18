---
layout: default
title: Dice
nav_order: 14
---

# Dice Rolling

`github.com/mechanical-lich/mlge/dice`

Parses and evaluates tabletop-style dice expressions like `2d6+3`, `1d20`, or `2d6+3-1d4`.

## Quick Roll

```go
import "github.com/mechanical-lich/mlge/dice"

result, err := dice.Roll("2d6")
// result is a random integer (2-12)
```

## Parsed Dice

For full expressions with modifiers, use `ParseDiceRequest`:

```go
d, err := dice.ParseDiceRequest("2d6+3-1d4")
fmt.Println(d.Result)    // Final computed value
fmt.Println(d.Breakdown) // Detailed roll breakdown
```

### Dice Struct

```go
type Dice struct {
    Result    int
    Breakdown string
}
```

| Field | Description |
|-------|-------------|
| `Result` | The final numeric result of the expression |
| `Breakdown` | A string showing each individual die roll |

## Expression Syntax

| Expression | Meaning |
|------------|---------|
| `1d6` | Roll one six-sided die |
| `2d6` | Roll two six-sided dice and sum |
| `2d6+3` | Roll 2d6 and add 3 |
| `2d6-1` | Roll 2d6 and subtract 1 |
| `2d6+1d4` | Roll 2d6 and add a 1d4 roll |
| `2d6+3-1d4` | Complex expression with multiple terms |

## Tokenizer

For custom parsing needs, `TokenizeDiceRequest` splits an expression into its component tokens:

```go
tokens := dice.TokenizeDiceRequest("2d6+3-1d4")
// ["2d6", "+", "3", "-", "1d4"]
```
