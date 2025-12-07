package balatro

import (
	"fmt"
	"golatro/pkg/balatro/rand"
	"strconv"
)

type GameState struct {
	Deck               []Card
	HandLevels         map[HandType]HandLevel
	MaxHands           int
	MaxDiscards        int
	MaxSelectableCards int
	MaxHandSize        int
	Jokers             []Joker
	MaxJokers          int
	Round              int
	Ante               int
	HighestScore       int
	Money              int
	MinMoney           int
	Consumables        []Consumable
	MaxConsumables     int
	GameLogger         *GameLogger

	HandTypeUsage map[HandType]int
	MaxRoundScore float64
	MaxHandScore  float64
}

func NewBasicGameState() GameState {
	initHandLevels := DefaultHandLevels()

	initCards := make([]Card, 0, 52)
	suitsList := []Suit{Clubs, Diamonds, Hearts, Spades}
	ranksList := []Rank{R2, R3, R4, R5, R6,
		R7, R8, R9, R10,
		J, Q, K, A}
	for _, s := range suitsList {
		for _, r := range ranksList {
			initCards = append(initCards, NewBasicCard(r, s))
		}
	}

	return GameState{
		Deck:               initCards,
		HandLevels:         initHandLevels,
		MaxHands:           4,
		MaxDiscards:        4,
		MaxSelectableCards: 8,
		MaxHandSize:        5,
		Jokers:             make([]Joker, 0),
		MaxJokers:          5,
		Round:              0,
		Ante:               1,
		HighestScore:       0,
		Money:              4,
		MinMoney:           0,
		Consumables:        make([]Consumable, 0),
		MaxConsumables:     0,
		GameLogger:         &GameLogger{log: make([]GameLog, 0)},

		HandTypeUsage: make(map[HandType]int),
		MaxRoundScore: 0,
		MaxHandScore:  0,
	}
}

const ROUNDS_PER_ANTE = 3
const WIN_ANTE = 8

func (g *GameState) NextRound() RoundState {
	availableCards := make([]Card, 0)
	availableCards = append(availableCards, g.Deck...)

	selectableCards := make([]Card, 0)
	for range g.MaxSelectableCards {
		randomIdx := rand.Int(len(availableCards))
		selectableCards = append(selectableCards, availableCards[randomIdx])
		if randomIdx == len(availableCards)-1 {
			availableCards = availableCards[:randomIdx]
		} else {
			availableCards = append(availableCards[:randomIdx], availableCards[randomIdx+1:]...)
		}
	}

	Hand(selectableCards).SortFunc(CompareRankSuit)

	if g.Round >= ROUNDS_PER_ANTE {
		g.Round = 1
		g.Ante += 1
	} else {
		g.Round += 1
	}
	g.GameLogger.Add(ROUND_STARTED, fmt.Sprintf("Ante %v, round %v", g.Ante, g.Round))

	return RoundState{
		AvailableCards:  availableCards,
		SelectableCards: selectableCards,
		Target:          TargetPoints(g.Round, g.Ante),
		Points:          0,
		Hand:            1,
		Discards:        0,
		logger:          g.GameLogger,
	}
}

func (g *GameState) TriggerJokers(t JokerTiming, sum int, mult float64, round *RoundState, hand Hand, cardIdx int, leftOverCards []Card) (int, float64) {
	for _, j := range g.Jokers {
		for _, e := range j.Type.Effects {
			if e.timing == t {
				s, m := e.Effect(g, round, hand, cardIdx, leftOverCards)
				if s != NoOpSum || m != NoOpMult {
					switch e.timing {
					case PerCard:
						g.GameLogger.Add(JOKER_ACTIVATED, fmt.Sprintf("%v. %v. %v", j, sumMultScoreDetails(s, m), hand[cardIdx]))
					case PerHand:
						g.GameLogger.Add(JOKER_ACTIVATED, fmt.Sprintf("%v. %v", j, sumMultScoreDetails(s, m)))
					case PerRound:
						g.GameLogger.Add(JOKER_ACTIVATED, fmt.Sprintf("%v", j))
					case BlindSelected:
						g.GameLogger.Add(JOKER_ACTIVATED, fmt.Sprintf("%v", j))
					}
				}
				sum += int(s)
				mult = ApplyMultiplier(mult, m)
			}
		}
	}

	return sum, mult
}

func (g *GameState) SellJoker(idx int, sellPrice func(Consumable) int) bool {
	if idx < 0 || idx >= len(g.Jokers) {
		return false
	}
	joker := g.Jokers[idx]
	price := sellPrice(JokerCard(joker))
	g.Money += price
	if idx == len(g.Jokers)-1 {
		g.Jokers = g.Jokers[:idx]
	} else {
		g.Jokers = append(g.Jokers[:idx], g.Jokers[idx+1:]...)
	}

	for _, e := range joker.Type.Effects {
		e.RemoveEffect(g, nil)
	}
	g.GameLogger.Add(JOKER_REMOVED, "Sold "+joker.Type.String()+" for "+strconv.Itoa(price)+"c")
	return true
}
