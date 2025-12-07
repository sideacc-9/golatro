package balatro

type Edition interface {
	Apply(int, Suit) (Sum, Multiplier)
	Detailer
	Abbreviater
}

var Editions = []Edition{
	Normal{},
}

type Normal struct{}

func (Normal) Apply(int, Suit) (Sum, Multiplier) {
	return NoOpSum, NoOpMult
}

func (p Normal) String() string {
	return ""
}

func (p Normal) Help() string {
	return ""
}

func (b Normal) Abbreviation() string {
	return ""
}
