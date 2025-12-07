package balatro

import "fmt"

type Detailer interface {
	fmt.Stringer
	Help() string
}

type Abbreviater interface {
	Abbreviation() string
}
