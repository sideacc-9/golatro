package balatro

import (
	"golatro/pkg/balatro/rand"

	"github.com/google/uuid"
)

type Rarity int

const (
	Common Rarity = iota
	Uncommon
	Rare
	Legendary
)

type JokerTiming int

const (
	PerCard JokerTiming = iota
	PerHand
	PerRound
	Passive
	BlindSelected
)

type JokerFunc func(game *GameState, round *RoundState, hand Hand, cardIdx int, leftOverCards []Card) (Sum, Multiplier)
type JokerPassive func(*GameState, *RoundState)

type JokerEffect struct {
	effect     JokerFunc
	undoEffect JokerPassive
	timing     JokerTiming
}

// plain joker functionality. is immutable and does not have enhacements
type JokerType struct {
	Effects []JokerEffect

	Rarity Rarity
	name   string
	help   string
}

func (j JokerEffect) Effect(game *GameState, round *RoundState, hand Hand, cardIdx int, leftOverCards []Card) (Sum, Multiplier) {
	if j.timing == PerCard && (cardIdx < 0 || cardIdx >= len(hand)) {
		panic("Effect is of timing per_card but received a non valid card index")
	}
	return j.effect(game, round, hand, cardIdx, leftOverCards)
}

func (j JokerEffect) RemoveEffect(game *GameState, round *RoundState) {
	if j.timing != Passive {
		return
	}
	j.undoEffect(game, round)
}

func NewEffect(timing JokerTiming, jFunc JokerFunc, jUndo JokerPassive) JokerEffect {
	return JokerEffect{
		effect:     jFunc,
		undoEffect: jUndo,
		timing:     timing,
	}
}

func NewBasicEffect(timing JokerTiming, jFunc JokerFunc) JokerEffect {
	if timing == Passive {
		panic("Use NewPassiveEffect to set undo effect")
	}
	return NewEffect(timing, jFunc, nil)
}

func NewPassiveEffect(jFunc JokerPassive, jUndo JokerPassive) JokerEffect {
	return NewEffect(Passive, func(game *GameState, round *RoundState, hand Hand, cardIdx int, leftOverCards []Card) (Sum, Multiplier) {
		jFunc(game, round)
		return NoOpSum, NoOpMult
	}, jUndo)
}

func (j JokerType) String() string {
	return j.name
}

func (j JokerType) Help() string {
	return j.help
}

func NewJokerType(name, help string, rarity Rarity, effects []JokerEffect) JokerType {
	return JokerType{
		Effects: effects,
		name:    name,
		help:    help,
		Rarity:  rarity,
	}
}

type Joker struct {
	Type        JokerType
	Edition     Edition
	Enhancement Enhancement
}

func NewJoker(t JokerType) Joker {
	return Joker{
		Type:        t,
		Edition:     Normal{},
		Enhancement: None{},
	}
}

var (
	ClassicJoker = NewJokerType("Joker", "Adds +4 to mult counter", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, _ Hand, _ int, _ []Card) (Sum, Multiplier) {
			return NoOpSum, Multiplier{Type: Add, Value: 4}
		})})

	GreedyJoker = NewJokerType("Greedy Joker", "Adds +3 for each "+Diamonds.Symbol+" suit in play", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			if hand[cardIdx].Suit == Diamonds {
				return NoOpSum, Multiplier{Type: Add, Value: 3}
			}
			return NoOpSum, NoOpMult
		})})
	LustyJoker = NewJokerType("Lusty Joker", "Adds +3 for each "+Hearts.Symbol+" suit in play", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			if hand[cardIdx].Suit == Hearts {
				return NoOpSum, Multiplier{Type: Add, Value: 3}
			}
			return NoOpSum, NoOpMult
		})})
	WrathfulJoker = NewJokerType("Wrathful Joker", "Adds +3 for each "+Spades.Symbol+" suit in play", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			if hand[cardIdx].Suit == Spades {
				return NoOpSum, Multiplier{Type: Add, Value: 3}
			}
			return NoOpSum, NoOpMult
		})})
	GluttonousJoker = NewJokerType("Gluttonous Joker", "Adds +3 for each "+Clubs.Symbol+" suit in play", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			if hand[cardIdx].Suit == Clubs {
				return NoOpSum, Multiplier{Type: Add, Value: 3}
			}
			return NoOpSum, NoOpMult
		})})

	JollyJoker = NewJokerType("Jolly Joker", "+8 Mult if played hand contains a Pair", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(Pair, hand); ok {
				return NoOpSum, Multiplier{Type: Add, Value: 8}
			}
			return NoOpSum, NoOpMult
		})})
	ZanyJoker = NewJokerType("Zany Joker", "+12 Mult if played hand contains a Three of a Kind", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(ThreeOfAKind, hand); ok {
				return NoOpSum, Multiplier{Type: Add, Value: 12}
			}
			return NoOpSum, NoOpMult
		})})
	MadJoker = NewJokerType("Mad Joker", "+10 Mult if played hand contains a Two Pair", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(TwoPair, hand); ok {
				return NoOpSum, Multiplier{Type: Add, Value: 10}
			}
			return NoOpSum, NoOpMult
		})})
	CrazyJoker = NewJokerType("Crazy Joker", "+12 Mult if played hand contains a Straight", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(Straight, hand); ok {
				return NoOpSum, Multiplier{Type: Add, Value: 12}
			}
			return NoOpSum, NoOpMult
		})})
	DrollJoker = NewJokerType("Droll Joker", "+10 Mult if played hand contains a Flush", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(Flush, hand); ok {
				return NoOpSum, Multiplier{Type: Add, Value: 10}
			}
			return NoOpSum, NoOpMult
		})})

	SlyJoker = NewJokerType("Sly Joker", "+50 Chips if played hand contains a Pair", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(Pair, hand); ok {
				return Sum(50), NoOpMult
			}
			return NoOpSum, NoOpMult
		})})
	WilyJoker = NewJokerType("Wily Joker", "+100 Chips if played hand contains a Three of a Kind", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(ThreeOfAKind, hand); ok {
				return Sum(100), NoOpMult
			}
			return NoOpSum, NoOpMult
		})})
	CleverJoker = NewJokerType("Clever Joker", "+80 Chips if played hand contains a Two Pair", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(TwoPair, hand); ok {
				return Sum(50), NoOpMult
			}
			return NoOpSum, NoOpMult
		})})
	DeviousJoker = NewJokerType("Devious Joker", "+100 Chips if played hand contains a Straight", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(Straight, hand); ok {
				return Sum(50), NoOpMult
			}
			return NoOpSum, NoOpMult
		})})
	CraftyJoker = NewJokerType("Crafty Joker", "+80 Chips if played hand contains a Flush", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if _, ok := CheckHandType(Flush, hand); ok {
				return Sum(50), NoOpMult
			}
			return NoOpSum, NoOpMult
		})})

	HalfJoker = NewJokerType("Half Joker", "+20 Mult if played hand contains 3 or fewer cards", Uncommon,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if len(hand) <= 3 {
				return NoOpSum, Multiplier{Type: Add, Value: 20}
			}
			return NoOpSum, NoOpMult
		})})

	FourFingers = NewJokerType("Four Fingers", "All Flushes and Straights can be made with 4 cards", Uncommon,
		[]JokerEffect{NewPassiveEffect(func(_ *GameState, _ *RoundState) {
			ChangeHandCondition(Flush, flushF(4))
			ChangeHandCondition(Straight, straightF(4))
		}, func(_ *GameState, _ *RoundState) {
			ChangeHandCondition(Flush, flushF(5))
			ChangeHandCondition(Straight, straightF(5))
		})})

	Mime = NewJokerType("Mime", "Retrigger all card held in hand abilities", Uncommon,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, hand Hand, cardIdx int, leftOverCards []Card) (Sum, Multiplier) {
			var sum, mult = NoOpSum, NoOpMult // not correct. It should start with the round's current mult and sum
			for _, c := range leftOverCards {
				if c.Enhancement.EffectTarget() == UNPLAYED {
					sumI, multI := c.Enhancement.Apply(hand[cardIdx].Rank, hand[cardIdx].Suit)
					sum += sumI
					mult = Multiplier{Type: Add, Value: ApplyMultiplier(mult.Value, multI)}
				}
			}
			return sum, mult
		})})

	CreditCard = NewJokerType("Credit Card", "Go up to -$20 in debt", Common,
		[]JokerEffect{NewPassiveEffect(func(game *GameState, _ *RoundState) {
			game.MinMoney = -20
		}, func(gs *GameState, _ *RoundState) {
			gs.MinMoney = 0
		})})

	Banner = NewJokerType("Banner", "+30 Chips for each remaining discard", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(game *GameState, round *RoundState, _ Hand, _ int, _ []Card) (Sum, Multiplier) {
			return Sum(30 * (game.MaxDiscards - round.Discards)), NoOpMult
		})})

	MysticSummit = NewJokerType("Mystic Summit", "+15 Mult when 0 discards remaining", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(game *GameState, round *RoundState, _ Hand, _ int, _ []Card) (Sum, Multiplier) {
			if game.MaxDiscards == round.Discards {
				return NoOpSum, Multiplier{Type: Add, Value: 15}
			}
			return NoOpSum, NoOpMult
		})})

	Misprint = NewJokerType("Misprint", "+0-23 Mult", Common,
		[]JokerEffect{NewBasicEffect(PerHand, func(_ *GameState, _ *RoundState, _ Hand, _ int, _ []Card) (Sum, Multiplier) {
			return NoOpSum, Multiplier{Value: float64(rand.Int(23)), Type: Add}
		})})

	// todo. Add current mult and chips to round state in order to use it here
	Dusk = NewJokerType("Dusk", "Retrigger all played cards in final hand of the round",
		Uncommon, []JokerEffect{NewBasicEffect(PerHand, func(game *GameState, round *RoundState, hand Hand, _ int, leftOverCards []Card) (Sum, Multiplier) {
			if round.Hand == game.MaxHands {
				sum, mult := scoreCards(game, round, leftOverCards, []Card(hand), 0, 1)
				return Sum(sum), Multiplier{Type: Add, Value: mult}
			}
			return NoOpSum, NoOpMult
		})})

	// todo Chaos the Clown. Add rerolls
	Fibonacci = NewJokerType("Fibonacci", "Each played Ace, 2, 3, 5, or 8 gives +8 Mult when scored",
		Uncommon, []JokerEffect{NewBasicEffect(PerCard, func(_ *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			r := hand[cardIdx].Rank
			if r == R2 || r == R3 || r == R5 || r == R8 {
				return NoOpSum, Multiplier{Type: Add, Value: +8}
			}
			return NoOpSum, NoOpMult
		})})
	SteelJoker = NewJokerType("Steel Joker", "Gives X0.2 Mult for each Steel Card in your full deck",
		Uncommon, []JokerEffect{NewBasicEffect(PerHand, func(game *GameState, _ *RoundState, hand Hand, cardIdx int, leftOverCards []Card) (Sum, Multiplier) {
			steel := NewSteel()
			nSteels := 0
			for _, c := range game.Deck {
				if c.Enhancement == steel {
					nSteels++
				}
			}
			return NoOpSum, Multiplier{Type: Multiply, Value: 1 + 0.2*float64(nSteels)}
		})})
	ScaryFace = NewJokerType("Scary Face", "Played face cards give +30 Chips when scored",
		Common, []JokerEffect{NewBasicEffect(PerCard, func(_ *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			r := hand[cardIdx].Rank
			if r == J || r == Q || r == K {
				return Sum(30), NoOpMult
			}
			return NoOpSum, NoOpMult
		})})
	AbstractJoker = NewJokerType("Abstract Joker", "+3 Mult for each Joker card",
		Common, []JokerEffect{NewBasicEffect(PerHand, func(game *GameState, _ *RoundState, _ Hand, _ int, _ []Card) (Sum, Multiplier) {
			return NoOpSum, Multiplier{Type: Add, Value: float64(len(game.Jokers))}
		})})

	// todo Delayed Gratification. I honestly dont even understand how the effect works

	// todo. same with dusk. Need current mult and sum to continue chain in scoreCards
	Hack = NewJokerType("Hack", "Retrigger each played 2, 3, 4, or 5",
		Common, []JokerEffect{NewBasicEffect(PerCard, func(game *GameState, round *RoundState, hand Hand, cardIdx int, leftOverCards []Card) (Sum, Multiplier) {
			r := hand[cardIdx].Rank
			if r != R2 && r != R3 && r != R4 && r != R5 {
				return NoOpSum, NoOpMult
			}
			card := hand[cardIdx]
			return card.Enhancement.Apply(card.Rank, card.Suit)
		})})

	// todo Pareidolia. Need some kind of isFace func to overwrite. Will think about it later

	EvenSteven = NewJokerType("Even Steven", "Played cards with even rank give +4 Mult when scored",
		Common, []JokerEffect{NewBasicEffect(PerCard, func(game *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			switch hand[cardIdx].Rank {
			case R2, R4, R6, R8, R10:
				return NoOpSum, Multiplier{Type: Add, Value: 4}
			default:
				return NoOpSum, NoOpMult
			}
		})})
	OddTodd = NewJokerType("Odd Todd", "Played cards with odd rank give +31 Chips when scored",
		Common, []JokerEffect{NewBasicEffect(PerCard, func(game *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			switch hand[cardIdx].Rank {
			case A, R3, R5, R7, R9:
				return Sum(31), NoOpMult
			default:
				return NoOpSum, NoOpMult
			}
		})})
	Scholar = NewJokerType("Scholar", "Played Aces give +20 Chips and +4 Mult when scored",
		Common, []JokerEffect{NewBasicEffect(PerCard, func(game *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			if hand[cardIdx].Rank == A {
				return Sum(20), Multiplier{Type: Add, Value: 4}
			}
			return NoOpSum, NoOpMult
		})})
	BusinessCard = NewJokerType("Business Card", "Played face cards have a 1 in 2 chance to give $2 when scored",
		Common, []JokerEffect{NewBasicEffect(PerCard, func(game *GameState, _ *RoundState, hand Hand, cardIdx int, _ []Card) (Sum, Multiplier) {
			if hand[cardIdx].Rank == A && rand.Int(2) == 0 {
				game.Money += 2
				game.GameLogger.Add(MONEY_GAINED, "Business Card added 2c")
			}
			return NoOpSum, NoOpMult
		})})

	// todo Supernova. Add num times played x hand

	// todo Ride the Bus. Add hand history to gamestate

	SpaceJoker = NewJokerType("Space Joker", "1 in 4 chance to upgrade level of played poker hand",
		Uncommon, []JokerEffect{NewBasicEffect(PerHand, func(game *GameState, _ *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if rand.Int(4) == 0 {
				for _, handType := range HandTypes {
					if _, ok := CheckHandType(handType, hand); ok {
						CelestialCard(handType).Apply(game)
						break
					}
				}
			}
			return NoOpSum, NoOpMult
		})})

	// todo Egg. Add individual cost/sell value to jokers

	Burglar = NewJokerType("Burglar", "Gain +3 Hands and lose all discards",
		Common, []JokerEffect{NewPassiveEffect(func(game *GameState, _ *RoundState) {
			game.MaxHands += 3
			game.MaxDiscards = 0
		}, func(game *GameState, _ *RoundState) {
			game.MaxHands -= 3
			game.MaxDiscards = 4
		})})

	DNA = NewJokerType("DNA", "If first hand of round has only 1 card, add a permanent copy to deck and draw it to hand",
		Rare, []JokerEffect{NewBasicEffect(PerHand, func(game *GameState, round *RoundState, hand Hand, _ int, _ []Card) (Sum, Multiplier) {
			if len(hand) == 1 && round.Hand == 1 {
				copy := hand[0]
				copy.Uuid = uuid.New()
				game.Deck = append(game.Deck, copy)
				round.AvailableCards = append(round.AvailableCards, copy)
				round.SelectableCards = append(round.SelectableCards, copy)
				game.GameLogger.Add(CARD_ADDED, "DNA duplicated "+copy.String())
			}
			return NoOpSum, NoOpMult
		})})
)

var JokerTypes = []JokerType{
	ClassicJoker,

	// +mult if card matches suit
	GreedyJoker,
	LustyJoker,
	WrathfulJoker,
	GluttonousJoker,

	// +mult if hand type
	JollyJoker,
	ZanyJoker,
	MadJoker,
	CrazyJoker,
	DrollJoker,

	// +chips if hand type
	SlyJoker,
	WilyJoker,
	CleverJoker,
	DeviousJoker,
	CraftyJoker,

	HalfJoker,
	FourFingers,
	// todo. Add current mult and chips to round state in order to use it here
	Mime,
	CreditCard,
	// todo Ceremonial Dagger. New joker timing, when blind selected
	Banner,
	MysticSummit,
	// todo Marble Joker. joker timing: blind selected
	// todo Loyalty Card. must have "hands played" counter
	// todo 8 Ball. implement tarot cards
	Misprint,
	// todo. Add current mult and chips to round state in order to use it here
	Dusk,
	// todo Chaos the Clown. Add rerolls
	Fibonacci,
	SteelJoker,
	ScaryFace,
	AbstractJoker,
	// todo Delayed Gratification. I honestly dont even understand how the effect works
	// todo. same with dusk. Need current mult and sum to continue chain in scoreCards
	Hack,
	// todo Pareidolia. Need some kind of isFace func to overwrite. Will think about it later
	EvenSteven,
	OddTodd,
	Scholar,
	BusinessCard,
	// todo Supernova. Add num times played x hand
	// todo Ride the Bus. Add hand history to gamestate
	SpaceJoker,
	// todo Egg. Add individual cost/sell value to jokers
	Burglar,
	DNA,
}
