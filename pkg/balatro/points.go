package balatro

import "math"

type MultType int

const (
	Add MultType = iota
	Multiply
	Exponent
)

func (m MultType) String() string {
	switch m {
	case Add:
		return "+"
	case Multiply:
		return "x"
	case Exponent:
		return "^"
	}
	return ""
}

type Multiplier struct {
	Value float64
	Type  MultType
}

var NoOpMult = Multiplier{Value: 0, Type: Add}

type Sum int

var NoOpSum = Sum(0)

func ApplyMultiplier(mult float64, m Multiplier) float64 {
	switch m.Type {
	case Add:
		mult += m.Value
	case Multiply:
		mult *= m.Value
	case Exponent:
		mult = math.Pow(mult, m.Value)
	}
	return mult
}
