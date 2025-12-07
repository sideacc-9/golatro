package balatro

import (
	"fmt"
	"golatro/pkg/balatro/rand"
)

type Equaler interface {
	Equals(any) bool
}

type Helper interface {
	Help() string
}

type Consumable interface {
	Apply(*GameState) error
	fmt.Stringer
	Price() int
	Equaler
	Helper
}

type Pack interface {
	// Get a set of cards to choose from and pack size for info about how many to choose from
	Open() []Consumable
	Size() PackSize
	fmt.Stringer
	Helper
}

type PackSize struct {
	probability    int
	Size           int
	ChoosableCards int
	Price          int
	Name           string
}

func (p PackSize) String() string {
	return p.Name
}

var (
	Mini = PackSize{
		probability:    3,
		Size:           3,
		ChoosableCards: 1,
		Price:          4,
		Name:           "Normal",
	}
	Jumbo = PackSize{
		probability:    2,
		Size:           5,
		ChoosableCards: 1,
		Price:          6,
		Name:           "Jumbo",
	}
	Mega = PackSize{
		probability:    1,
		Size:           5,
		ChoosableCards: 2,
		Price:          8,
		Name:           "Mega",
	}
)

func (size PackSize) Help() string {
	return fmt.Sprintf("%v, %v cards, choose %v", size.Name, size.Size, size.ChoosableCards)
}

var sizes = []PackSize{Mini, Jumbo, Mega}
var probabilityPool = make([]int, 0)

func init() {
	for idx, size := range sizes {
		for i := 0; i < size.probability; i++ {
			probabilityPool = append(probabilityPool, idx)
		}
	}
}

func getRandomPackSize() PackSize {
	packSizeIndex := rand.Int(len(probabilityPool))
	return sizes[probabilityPool[packSizeIndex]]
}

func GetRandomPack() Pack {
	packSize := getRandomPackSize()
	switch rand.Int(3) {
	case 0:
		return CelestialPack{
			size: packSize,
		}
	case 1:
		return PlayingCardPack{
			size: packSize,
		}
	// case 2:	// todo: arcana and spectral packs
	default:
		return JokerCardPack{
			size: packSize,
		}
	}
}

func GetRandomConsumable() Consumable {
	switch rand.Int(3) {
	case 0:
		return GetRandomHandType()
	case 1:
		return GetRandomCard()
	// case 2:	// todo: arcana and spectral packs
	default:
		return GetRandomJoker()
	}
}

type CelestialCard HandType

func (c CelestialCard) String() string {
	return HandType(c).String()
}

func (c CelestialCard) Help() string {
	return "Upgrade hand type '" + c.String() + "'"
}

func (c CelestialCard) Apply(g *GameState) error {
	handlevel := g.HandLevels[HandType(c)]
	handlevel.Upgrade()
	g.HandLevels[HandType(c)] = handlevel
	return nil
}

func (c CelestialCard) Price() int {
	return 3
}

func (c CelestialCard) Equals(other any) bool {
	otherCelestial, ok := other.(CelestialCard)
	if ok {
		return otherCelestial == c
	}
	return false
}

type CelestialPack struct {
	size PackSize
}

func GetRandomHandType() CelestialCard {
	randIdx := rand.Int(len(HandTypes))
	card := HandTypes[randIdx]
	return CelestialCard(card)
}

func (p CelestialPack) Open() []Consumable {
	packSize := p.size
	cards := make([]Consumable, 0, packSize.Size)
	for i := 0; i < packSize.Size; i++ {
		cards = append(cards, GetRandomHandType())
	}
	return cards
}

func (p CelestialPack) Size() PackSize {
	return p.size
}

func (p CelestialPack) String() string {
	return "Celestial"
}

func (p CelestialPack) Help() string {
	return "Choose upgrades for hand levels\n\t(" + p.size.Help() + ")"
}

type PlayingCard Card

func (p PlayingCard) String() string {
	return fmt.Sprintf("%v %v\n%v", p.Suit.Symbol, p.Rank, p.Enhancement)
}

func (c PlayingCard) Apply(g *GameState) error {
	g.Deck = append(g.Deck, Card(c))
	g.GameLogger.Add(CARD_ADDED, "Added "+Card(c).String())
	return nil
}

func (p PlayingCard) Help() string {
	str := "Add a " + p.Suit.Symbol + " " + p.Rank.String()
	return str
}

func (c PlayingCard) Price() int {
	card := Card(c)
	price := 2
	normalEd := Normal{}
	if card.Edition != normalEd {
		price += 2
	}
	noneEnh := None{}
	if card.Enhancement != noneEnh {
		price += 2
	}
	return price
}

func (c PlayingCard) Equals(other any) bool {
	otherCard, ok := other.(PlayingCard)
	if ok {
		return otherCard == c
	}
	return false
}

type PlayingCardPack struct {
	size PackSize
}

func randEdition(chancePercentNormal int) Edition {
	isNormal := rand.Int(100)
	if isNormal < chancePercentNormal {
		return Normal{}
	}
	randIdx := rand.Int(len(Editions))
	return Editions[randIdx]
}

func randEnhancement(chancePercentNone int) Enhancement {
	isNormal := rand.Int(100)
	if isNormal < chancePercentNone {
		return None{}
	}
	randIdx := rand.Int(len(Enhancements))
	return Enhancements[randIdx]
}

var none = None{}

func GetRandomCard() PlayingCard {
	card := NewBasicCard(RandomRank(), RandomSuit())

	card.Enhancement = randEnhancement(40)
	chancesNormal := 40
	if card.Enhancement != none { // reduce chances to get special edition if the card has an enhancement
		chancesNormal += 20
	}
	card.Edition = randEdition(chancesNormal)

	return PlayingCard(card)
}

func (p PlayingCardPack) Open() []Consumable {
	packSize := p.size
	cards := make([]Consumable, 0, packSize.Size)
	for i := 0; i < packSize.Size; i++ {
		cards = append(cards, GetRandomCard())
	}
	return cards
}

func (p PlayingCardPack) Size() PackSize {
	return p.size
}

func (p PlayingCardPack) String() string {
	return "Standard"
}

func (p PlayingCardPack) Help() string {
	return "Choose a card to add to your deck\n\t(" + p.size.Help() + ")"
}

type JokerCard Joker

func (c JokerCard) Apply(g *GameState) error {
	if len(g.Jokers) >= g.MaxJokers {
		return fmt.Errorf("all joker slots occupied")
	}
	g.Jokers = append(g.Jokers, Joker(c))

	for _, e := range c.Type.Effects {
		if e.timing == Passive {
			e.Effect(g, nil, nil, -1, nil)
		}
	}

	return nil
}

func (c JokerCard) String() string {
	return JokerType(c.Type).String()
}

func (c JokerCard) Help() string {
	return JokerType(c.Type).Help()
}

func (c JokerCard) Price() int {
	j := JokerType(c.Type)
	price := (int(j.Rarity) + 1) * 3
	normalEd := Normal{}
	if c.Edition != normalEd {
		price += 3
	}
	noneEnh := None{}
	if c.Enhancement != noneEnh {
		price += 3
	}
	return price
}

func (c JokerCard) Equals(other any) bool {
	otherJoker, ok := other.(JokerCard)
	if ok {
		return otherJoker.Type.name == c.Type.name &&
			otherJoker.Edition == c.Edition &&
			otherJoker.Enhancement == c.Enhancement
	}
	return false
}

type JokerCardPack struct {
	size PackSize
}

func GetRandomJoker() JokerCard {
	randIdx := rand.Int(len(JokerTypes))
	joker := JokerTypes[randIdx]
	return JokerCard(NewJoker(joker))
}

func (p JokerCardPack) Open() []Consumable {
	packSize := p.size
	cards := make([]Consumable, 0, packSize.Size)
	for i := 0; i < packSize.Size; i++ {
		cards = append(cards, GetRandomJoker())
	}
	return cards
}

func (p JokerCardPack) String() string {
	return "Buffoon"
}

func (p JokerCardPack) Size() PackSize {
	return p.size
}

func (p JokerCardPack) Help() string {
	return "Choose a joker to add\n\t(" + p.size.Help() + ")"
}
