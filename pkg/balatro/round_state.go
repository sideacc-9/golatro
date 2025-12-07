package balatro

import (
	"cmp"
	"fmt"
	"golatro/pkg/balatro/rand"
	"math"
	"slices"
	"strconv"
)

type RoundState struct {
	AvailableCards  []Card
	SelectableCards []Card
	Target          float64
	Points          float64
	Hand            int
	Discards        int
	logger          *GameLogger
}

type GameStatus int

const (
	WON = iota
	LOST
	INPROGRESS
)

func (r RoundState) RoundStatus(maxHands int) GameStatus {
	if r.Points >= r.Target {
		return WON
	}
	if r.Hand > maxHands {
		return LOST
	}
	return INPROGRESS
}

func (r *RoundState) Discard(hand Hand, maxSelectable int) {
	r.logger.Add(HAND_DISCARDED, fmt.Sprint(hand))
	r.removePlayedHand(hand, maxSelectable)
	r.Discards += 1
}

func (r *RoundState) removePlayedHand(hand Hand, maxSelectable int) {
	r.SelectableCards = slices.DeleteFunc(r.SelectableCards, func(c Card) bool {
		return slices.Contains(hand, c)
	})
	for i := len(r.SelectableCards); i < maxSelectable; i++ {
		randomIdx := rand.Int(len(r.AvailableCards))
		r.SelectableCards = append(r.SelectableCards, r.AvailableCards[randomIdx])
		if randomIdx == len(r.AvailableCards)-1 {
			r.AvailableCards = r.AvailableCards[:randomIdx]
		} else {
			r.AvailableCards = append(r.AvailableCards[:randomIdx], r.AvailableCards[randomIdx+1:]...)
		}
	}
	r.AvailableCards = slices.DeleteFunc(r.AvailableCards, func(c Card) bool {
		return slices.Contains(r.SelectableCards, c)
	})
}

func startAnte(ante float64) float64 {
	baseDiffAnte := 300.0
	return baseDiffAnte*math.Pow(2, ante) - baseDiffAnte
}

func TargetPoints(round, ante int) float64 {
	increment := 150.0
	anteF := float64(ante)

	start := startAnte(anteF)

	result := start + increment*math.Pow(2, anteF-1)*float64(round-1)
	const multipleOf float64 = 50
	return result / multipleOf * multipleOf
}

func endRound(game *GameState, round *RoundState) {
	if round.RoundStatus(game.MaxHands) == WON {
		initMoney := game.Money
		game.Money += 3
		game.Money += game.MaxHands - round.Hand
		game.Money += int(math.Min(5, float64(initMoney/5)))
		game.GameLogger.Add(MONEY_GAINED, fmt.Sprintf("Gained money: %v", game.Money-initMoney))
	}

	game.TriggerJokers(PerRound, 0, 0, round, Hand{}, -1, nil)

	if round.Points > game.MaxHandScore {
		game.MaxRoundScore = round.Points
	}
}

func filterStone(h Hand) (Hand, []Card) {
	stoneLessHand := make(Hand, 0)
	stones := make([]Card, 0)
	stoneEnh := NewStone()
	for _, c := range h {
		if c.Enhancement == stoneEnh {
			stones = append(stones, c)
		} else {
			stoneLessHand = append(stoneLessHand, c)
		}
	}
	return stoneLessHand, stones
}

// Applies scoring to a played hand according to current game and round states. Updates the round by removing the cards played
// from selectable cards, and fills the list with cards from available cards.
// Stone cards are treated separately. They cannot be used for hand type check but always count towards hand points
// Order of application:
//  1. start with current hand level sum and mult
//  2. run through cards that apply to handtype and apply all PerCard jokers to all of them
//  3. apply jokers PerHand jokers once
//  4. apply unscored card enhancements (steel cards)
func (round *RoundState) ScoreHand(game *GameState, hand Hand) (*HandType, int, float64) {
	defer endRound(game, round)

	game.GameLogger.Add(HAND_PLAYED, fmt.Sprint(hand))
	if len(hand) == 0 {
		return nil, 0, 0
	}
	stoneLess, stones := filterStone(hand)
	for _, handType := range HandTypes {
		if validCards, ok := CheckHandType(handType, stoneLess); ok {
			validCards = append(validCards, stones...)
			unscoredCards := filterUnscored(round.SelectableCards, []Card(hand))

			hLevel := game.HandLevels[handType]
			initSum, initMult := hLevel.GetSumMult()
			sum, mult := initSum, initMult
			sum, mult = scoreCards(game, round, unscoredCards, validCards, sum, mult)
			sum, mult = game.TriggerJokers(PerHand, sum, mult, round, hand, -1, unscoredCards)
			sum, mult = applyUnscoredCardEnhancements(unscoredCards, sum, mult)

			round.removePlayedHand(hand, game.MaxSelectableCards)
			round.Hand += 1
			round.Points += float64(sum) * mult

			if float64(sum)*mult > game.MaxHandScore {
				game.MaxHandScore = round.Points
			}

			checkBreakGlassCards(game, validCards)

			game.GameLogger.Add(HAND_SCORED, fmt.Sprintf("%v (lvl %v, initial %vx%v) - [%vx%v]: %v", handType, hLevel.Level,
				initSum, initMult, sum, mult, Hand(validCards)))
			return &handType, sum, mult
		}
	}
	return nil, 0, 0
}

func checkBreakGlassCards(game *GameState, playedHand Hand) {
	glass := NewGlass()
	for _, c := range playedHand {
		if c.Enhancement == glass && rand.Int(100) < glass.chanceBreak {
			game.Deck = slices.DeleteFunc(game.Deck, func(card Card) bool {
				return card.Uuid == c.Uuid
			})
			game.GameLogger.Add(CARD_DESTROYED, "Destroyed "+c.String())
		}
	}
}

func filterUnscored(total, played []Card) []Card {
	unscored := make([]Card, 0, len(total)-len(played))
	for _, c := range total {
		isPlayed := slices.ContainsFunc(played, func(card Card) bool {
			return cmp.Compare(card.Uuid.ID(), c.Uuid.ID()) == 0
		})
		if !isPlayed {
			unscored = append(unscored, c)
		}
	}
	return unscored
}

// Count score for all cards using the specified initial sum and multiplier
func scoreCards(game *GameState, round *RoundState, unscoredCards []Card, cards Hand, sum int, mult float64) (int, float64) {
	for i, card := range cards {
		var s Sum
		var m Multiplier
		if card.Enhancement.EffectTarget() == PLAYED {
			s, m = card.Enhancement.Apply(card.Rank, card.Suit)
			sum += int(s)
			mult = ApplyMultiplier(mult, m)
		} else {
			s, m = None{}.Apply(card.Rank, card.Suit)
			sum += int(s)
			mult = ApplyMultiplier(mult, m)
		}

		sum, mult = game.TriggerJokers(PerCard, sum, mult, round, cards, i, unscoredCards)

	}
	return sum, mult
}

func applyUnscoredCardEnhancements(cards []Card, sum int, mult float64) (int, float64) {
	for _, c := range cards {
		if c.Enhancement.EffectTarget() == UNPLAYED {
			s, m := c.Enhancement.Apply(c.Rank, c.Suit)
			sum += int(s)
			mult = ApplyMultiplier(mult, m)
		}
	}
	return int(sum), mult
}

func sumMultScoreDetails(s Sum, m Multiplier) string {
	pointsDetail := ""
	if s != NoOpSum {
		pointsDetail += "Sum: +" + strconv.Itoa(int(s))
	}
	if s != NoOpSum && m != NoOpMult {
		pointsDetail += ", "
	}
	if m != NoOpMult {
		pointsDetail += fmt.Sprintf("Mult: %v%v", m.Type.String(), m.Value)
	}
	return pointsDetail
}
