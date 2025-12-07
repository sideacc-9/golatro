package balatro

import (
	"fmt"
	"strconv"
)

type EffectTarget int

const (
	PLAYED EffectTarget = iota
	UNPLAYED
	DISCARDED
)

type Enhancement interface {
	Apply(Rank, Suit) (Sum, Multiplier)
	EffectTarget() EffectTarget
	Detailer
	Abbreviater
}

var Enhancements = []Enhancement{
	NewBonus(), NewMult(), NewSteel(), NewStone(), NewGlass(), None{},
}

type None struct{}

func (b None) Apply(rank Rank, _ Suit) (Sum, Multiplier) {
	return Sum(RankValue(rank)), NoOpMult
}

func (b None) EffectTarget() EffectTarget {
	return PLAYED
}

func (b None) String() string {
	return ""
}

func (b None) Help() string {
	return ""
}

func (b None) Abbreviation() string {
	return ""
}

func NewBonus() Bonus {
	return Bonus{Value: 30}
}

type Bonus struct {
	Value int
}

func (b Bonus) Apply(rank Rank, _ Suit) (Sum, Multiplier) {
	return Sum(RankValue(rank) + b.Value), NoOpMult
}

func (b Bonus) EffectTarget() EffectTarget {
	return PLAYED
}

func (b Bonus) String() string {
	return "Bonus"
}

func (b Bonus) Help() string {
	return fmt.Sprintf("Adds an extra %v towards the sum counter", b.Value)
}

func (b Bonus) Abbreviation() string {
	return "Bns"
}

func NewMult() Mult {
	return Mult{Value: 4}
}

type Mult struct {
	Value float64
}

func (b Mult) Apply(r Rank, _ Suit) (Sum, Multiplier) {
	return Sum(RankValue(r)), Multiplier{Value: b.Value, Type: Add}
}

func (b Mult) EffectTarget() EffectTarget {
	return PLAYED
}

func (b Mult) String() string {
	return "Mult"
}

func (b Mult) Help() string {
	return fmt.Sprintf("Adds an extra %v towards the mult counter", b.Value)
}

func (b Mult) Abbreviation() string {
	return "Mlt"
}

func NewSteel() Steel {
	return Steel{Value: 1.5}
}

type Steel struct {
	Value float64
}

func (b Steel) Apply(_ Rank, _ Suit) (Sum, Multiplier) {
	return NoOpSum, Multiplier{Value: b.Value, Type: Multiply}
}

func (b Steel) EffectTarget() EffectTarget {
	return UNPLAYED
}

func (b Steel) String() string {
	return "Steel"
}

func (b Steel) Help() string {
	return fmt.Sprintf("Muliplies mult counter by %v if not used in hand", b.Value)
}

func (b Steel) Abbreviation() string {
	return "Stl"
}

func NewStone() Stone {
	return Stone(50)
}

type Stone int

func (b Stone) Apply(_ Rank, _ Suit) (Sum, Multiplier) {
	return Sum(b), NoOpMult
}

func (b Stone) EffectTarget() EffectTarget {
	return PLAYED
}

func (b Stone) String() string {
	return "Stone"
}

func (b Stone) Help() string {
	return fmt.Sprintf("Counts as +%v", b)
}

func (b Stone) Abbreviation() string {
	return "Stn"
}

func NewGlass() Glass {
	return Glass{
		chanceBreak: 25,
		mult:        2,
	}
}

type Glass struct {
	chanceBreak int // 100/chanceBreak = %
	mult        int
}

func (g Glass) Apply(r Rank, _ Suit) (Sum, Multiplier) {
	return Sum(r.Value()), Multiplier{Value: float64(g.mult), Type: Multiply}
}

func (g Glass) EffectTarget() EffectTarget {
	return PLAYED
}

func (g Glass) String() string {
	return "Glass"
}

func (g Glass) Help() string {
	return strconv.Itoa(g.mult) + "x round multiplier. " + strconv.Itoa(g.chanceBreak) + "% chance to break"
}

func (g Glass) Abbreviation() string {
	return "Gls"
}
