package balatro

import (
	"cmp"
	"golatro/pkg/balatro/rand"
	"strconv"

	"github.com/google/uuid"
)

type Suit struct {
	Id     int
	Symbol string
	Name   string
}

var (
	Hearts   = Suit{Id: 0, Name: "Hearts", Symbol: "♥"}
	Clubs    = Suit{Id: 1, Name: "Clubs", Symbol: "♣"}
	Diamonds = Suit{Id: 2, Name: "Diamonds", Symbol: "♦"}
	Spades   = Suit{Id: 3, Name: "Spades", Symbol: "♠"}
	Any      = Suit{Id: 4, Name: "Wildcard", Symbol: "\u20DD "}
)

func (s Suit) String() string {
	return s.Symbol
}

var suitList = []Suit{Hearts, Clubs, Diamonds, Spades}

func RandomSuit() Suit {
	return suitList[rand.Int(len(suitList))]
}

type Rank int

const (
	R2 Rank = iota + 2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	J
	Q
	K
	A
)

func (r Rank) Value() int {
	return RankValue(r)
}

func RankValue(r Rank) int {
	switch r {
	case R2, R3, R4, R5, R6, R7, R8, R9, R10:
		return int(r)
	case J, Q, K:
		return 10
	case A:
		return 11
	}
	return 0
}

func (r Rank) String() string {
	switch r {
	case R2, R3, R4, R5, R6, R7, R8, R9, R10:
		return strconv.Itoa(int(r))
	case J:
		return "J"
	case Q:
		return "Q"
	case K:
		return "K"
	case A:
		return "A"
	}
	return "err"
}

var ranksList = []Rank{R2, R3, R4, R5, R6, R7, R8, R9, R10, J, Q, K, A}

func RandomRank() Rank {
	return ranksList[rand.Int(len(ranksList))]
}

type Card struct {
	Uuid        uuid.UUID
	Rank        Rank
	Suit        Suit
	Edition     Edition
	Enhancement Enhancement
}

func NewCard(r Rank, s Suit, edition Edition, enhancement Enhancement) Card {
	return Card{
		Uuid:        uuid.New(),
		Rank:        r,
		Suit:        s,
		Edition:     edition,
		Enhancement: enhancement,
	}
}

func NewBasicCard(r Rank, s Suit) Card {
	return Card{
		Uuid:        uuid.New(),
		Rank:        r,
		Suit:        s,
		Edition:     Normal{},
		Enhancement: None{},
	}
}

func CompareRankSuit(c1, c2 Card) int {
	if c1.Rank == c2.Rank {
		return cmp.Compare(c1.Suit.Id, c2.Suit.Id)
	}
	return cmp.Compare(c1.Rank, c2.Rank)
}

func CompareSuitRank(c1, c2 Card) int {
	if c1.Suit.Id == c2.Suit.Id {
		return cmp.Compare(c1.Rank, c2.Rank)
	}
	return cmp.Compare(c1.Suit.Id, c2.Suit.Id)
}

func (c Card) String() string {
	str := c.Suit.Symbol + " " + c.Rank.String()
	edStr := c.Edition.String()
	ehStr := c.Enhancement.String()

	if edStr != "" || ehStr != "" {
		str += " ("
		if edStr != "" {
			str += edStr
		}
		if edStr != "" && ehStr != "" {
			str += " "
		}
		if ehStr != "" {
			str += ehStr
		}
		str += ")"
	}
	return str
}

type SortType int

const (
	RANK SortType = iota
	SUIT
)
