package balatro

import (
	"bytes"
	"fmt"
	"slices"
)

type Hand []Card

func (h Hand) Sort() {
	slices.SortFunc(h, CompareRankSuit)
}

func (h Hand) SortFunc(compare func(c1 Card, c2 Card) int) {
	slices.SortFunc(h, compare)
}

func (h Hand) Sorted() Hand {
	tmp := Hand(make([]Card, len(h)))
	copy(tmp, h)
	tmp.Sort()
	return tmp
}

func (h Hand) SortedFunc(compare func(c1 Card, c2 Card) int) Hand {
	tmp := Hand(make([]Card, len(h)))
	copy(tmp, h)
	tmp.SortFunc(compare)
	return tmp
}

func (h Hand) String() string {
	buf := new(bytes.Buffer)
	for i, c := range h {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprint(buf, c)
	}
	return buf.String()
}

type HandType int

const (
	FiveOfAKind HandType = iota
	RoyalFlush
	StraightFlush
	FourOfAKind
	FullHouse
	Flush
	Straight
	ThreeOfAKind
	TwoPair
	Pair
	HighCard
)

func (h HandType) String() string {
	switch h {
	case FiveOfAKind:
		return "Five of a Kind"
	case RoyalFlush:
		return "Royal Flush"
	case StraightFlush:
		return "Straight Flush"
	case FourOfAKind:
		return "Four of a Kind"
	case FullHouse:
		return "Full House"
	case Flush:
		return "Flush"
	case Straight:
		return "Straight"
	case ThreeOfAKind:
		return "Three of a Kind"
	case TwoPair:
		return "Two Pair"
	case Pair:
		return "Pair"
	case HighCard:
		return "High Card"
	default:
		return "Unknown"
	}
}

type HandCondition func(Hand) ([]Card, bool)

var HandTypes = []HandType{RoyalFlush, StraightFlush, FourOfAKind, FullHouse, Flush, Straight, ThreeOfAKind, TwoPair, Pair, HighCard}

func straightF(n int) func(c Hand) ([]Card, bool) {
	return func(c Hand) ([]Card, bool) {
		if len(c) < n {
			return nil, false
		}
		cSorted := c.SortedFunc(CompareRankSuit)
		for i := 1; i < n; i++ {
			if cSorted[i].Rank != cSorted[i-1].Rank+1 {
				return nil, false
			}
		}
		return c, true
	}
}
func flushF(n int) func(c Hand) ([]Card, bool) {
	return func(c Hand) ([]Card, bool) {
		if len(c) < n {
			return nil, false
		}
		matches := make(map[Suit]int)
		for _, card := range c {
			matches[card.Suit]++
		}
		for _, v := range matches {
			if v >= n {
				return c, true
			}
		}
		return nil, false
	}
}

var handConditions = map[HandType]HandCondition{
	HighCard: func(c Hand) ([]Card, bool) {
		if len(c) == 0 {
			return nil, false
		}
		cSorted := c.SortedFunc(CompareRankSuit)

		return cSorted[len(c)-1:], true
	},
	Pair: func(c Hand) ([]Card, bool) {
		if len(c) < 2 {
			return nil, false
		}
		rankChecked := make(map[Rank]Card)
		for _, card := range c {
			if otherCard, ok := rankChecked[card.Rank]; ok {
				return []Card{otherCard, card}, true
			}
			rankChecked[card.Rank] = card
		}
		return nil, false
	},
	ThreeOfAKind: func(c Hand) ([]Card, bool) {
		if len(c) < 3 {
			return nil, false
		}
		rankChecked := make(map[Rank][]Card)
		for _, card := range c {
			if _, ok := rankChecked[card.Rank]; ok {
				rankChecked[card.Rank] = append(rankChecked[card.Rank], card)
				if len(rankChecked[card.Rank]) >= 3 {
					return rankChecked[card.Rank], true
				}
			} else {
				rankChecked[card.Rank] = []Card{card}
			}
		}
		return nil, false
	},
	FourOfAKind: func(c Hand) ([]Card, bool) {
		if len(c) < 4 {
			return nil, false
		}
		rankChecked := make(map[Rank][]Card)
		for _, card := range c {
			if _, ok := rankChecked[card.Rank]; ok {
				rankChecked[card.Rank] = append(rankChecked[card.Rank], card)
				if len(rankChecked[card.Rank]) >= 4 {
					return rankChecked[card.Rank], true
				}
			} else {
				rankChecked[card.Rank] = []Card{card}
			}
		}
		return nil, false
	},
	FiveOfAKind: func(c Hand) ([]Card, bool) {
		if len(c) < 5 {
			return nil, false
		}
		rankChecked := make(map[Rank][]Card)
		for _, card := range c {
			if _, ok := rankChecked[card.Rank]; ok {
				rankChecked[card.Rank] = append(rankChecked[card.Rank], card)
				if len(rankChecked[card.Rank]) >= 5 {
					return rankChecked[card.Rank], true
				}
			} else {
				rankChecked[card.Rank] = make([]Card, 0)
			}
		}
		return nil, false
	},
	TwoPair: func(c Hand) ([]Card, bool) {
		if len(c) < 4 {
			return nil, false
		}
		rankChecked := make(map[Rank][]Card)
		for _, card := range c {
			if instances, ok := rankChecked[card.Rank]; ok {
				rankChecked[card.Rank] = append(instances, card)
			} else {
				rankChecked[card.Rank] = []Card{card}
			}
		}
		pairs := [][]Card{}
		for _, instances := range rankChecked {
			if len(instances) == 2 {
				pairs = append(pairs, instances)
			}
		}
		if len(pairs) == 2 {
			return append(pairs[0], pairs[1]...), true
		}
		return nil, false
	},
	Flush: flushF(5),
	FullHouse: func(c Hand) ([]Card, bool) {
		if len(c) < 5 {
			return nil, false
		}
		rank1, rank2 := Rank(-1), Rank(-1)
		instances1, instances2 := make([]Card, 0), make([]Card, 0)
		for _, card := range c {
			if rank1 == Rank(-1) {
				rank1 = card.Rank
			} else if rank2 == Rank(-1) && card.Rank != rank1 {
				rank2 = card.Rank
			}

			if rank1 == card.Rank {
				instances1 = append(instances1, card)
			}
			if rank2 == card.Rank {
				instances2 = append(instances2, card)
			}
		}
		if (len(instances1) >= 2 && len(instances2) >= 3) || (len(instances2) >= 2 && len(instances1) >= 3) {
			return c, true
		}
		return nil, false
	},
	Straight: straightF(5),
}

func init() {
	handConditions[StraightFlush] = func(c Hand) ([]Card, bool) {
		cardsStraight, isStraight := handConditions[Straight](c)
		if !isStraight {
			return nil, false
		}
		_, isFlush := handConditions[Flush](cardsStraight)
		if !isFlush {
			return nil, false
		}
		return cardsStraight, true
	}

	handConditions[RoyalFlush] = func(c Hand) ([]Card, bool) {
		cardsSF, isSF := handConditions[StraightFlush](c)
		if !isSF {
			return nil, false
		}
		if cardsSF[0].Rank == A && cardsSF[4].Rank == R10 {
			return cardsSF, true
		}
		return nil, false
	}
}

func CheckHandType(hType HandType, hand Hand) ([]Card, bool) {
	return handConditions[hType](hand)
}

func ChangeHandCondition(hType HandType, condition HandCondition) {
	handConditions[hType] = condition
}

type HandLevel struct {
	Level       int
	initSum     int
	initMult    float64
	sumUpgrade  int
	multUpgrade float64
}

func newHandLevel(initSum, sumUpgrade int, initMult, multUpgrade float64) HandLevel {
	return HandLevel{
		Level:       1,
		initSum:     initSum,
		initMult:    initMult,
		sumUpgrade:  sumUpgrade,
		multUpgrade: multUpgrade,
	}
}

func (hl HandLevel) GetSumMult() (int, float64) {
	return hl.initSum + (hl.Level-1)*hl.sumUpgrade, hl.initMult + float64(hl.Level-1)*hl.multUpgrade
}

func (hl *HandLevel) Upgrade() {
	hl.Level += 1
}

func DefaultHandLevels() map[HandType]HandLevel {
	return map[HandType]HandLevel{
		FiveOfAKind:   newHandLevel(120, 35, 12, 3),
		RoyalFlush:    newHandLevel(110, 30, 10, 3),
		StraightFlush: newHandLevel(100, 40, 8, 3),
		FourOfAKind:   newHandLevel(60, 30, 7, 3),
		FullHouse:     newHandLevel(40, 25, 4, 2),
		Flush:         newHandLevel(35, 15, 4, 2),
		Straight:      newHandLevel(30, 30, 4, 2),
		ThreeOfAKind:  newHandLevel(30, 20, 3, 2),
		TwoPair:       newHandLevel(20, 20, 2, 1),
		Pair:          newHandLevel(10, 15, 2, 1),
		HighCard:      newHandLevel(5, 10, 1, 1),
	}
}
