package menu

import (
	"golatro/internal/balatro/menu"
	"golatro/pkg/balatro"
	"testing"
)

func TestDisplayPackCards(t *testing.T) {
	upgradesRaw := balatro.GetRandomPack().Open()
	upgrades := make([]menu.PackCardOption, 0, len(upgradesRaw))
	for i, u := range upgradesRaw {
		upgrades = append(upgrades, menu.PackCardOption{Card: u, Id: i})
	}
	menu.DisplayPackCards(upgrades, []menu.PackCardOption{upgrades[0]}, 0)
}
