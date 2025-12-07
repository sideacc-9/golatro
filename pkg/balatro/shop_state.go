package balatro

type ShopState struct {
	Packs       []Pack
	Upgrades    []Consumable
	RerollPrice int
}

func NewShopState(packsN, upgradesN int) ShopState {
	packs := make([]Pack, 0, packsN)
	upgrades := make([]Consumable, 0, upgradesN)
	return ShopState{
		Packs:       packs,
		Upgrades:    upgrades,
		RerollPrice: 5,
	}
}

func (shop *ShopState) RemoveUpgrade(c Consumable) {
	idx := -1
	for i, up := range shop.Upgrades {
		if up.Equals(c) {
			idx = i
		}
	}
	var rightSide []Consumable
	if idx == len(shop.Upgrades)-1 {
		rightSide = []Consumable{}
	} else {
		rightSide = shop.Upgrades[idx+1:]
	}
	shop.Upgrades = append(shop.Upgrades[:idx], rightSide...)
}

func (shop *ShopState) RemovePack(p Pack) {
	idx := -1
	for i, pack := range shop.Packs {
		if pack == p {
			idx = i
		}
	}
	var rightSide []Pack
	if idx == len(shop.Packs)-1 {
		rightSide = []Pack{}
	} else {
		rightSide = shop.Packs[idx+1:]
	}
	shop.Packs = append(shop.Packs[:idx], rightSide...)
}
