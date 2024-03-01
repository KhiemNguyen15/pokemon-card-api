package card

type Card struct {
	ID       int
	Name     string
	Number   string
	Rarity   string
	Value    float64
	ImageURL string
	Set      Set
}

type Set struct {
	Name   string
	Series string
	Total  int
}
