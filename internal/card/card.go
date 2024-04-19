package card

type Card struct {
	ID       int     `json:"id"        db:"id"`
	Name     string  `json:"name"      db:"name"`
	Number   string  `json:"number"    db:"number"`
	Rarity   string  `json:"rarity"    db:"rarity"`
	Value    float64 `json:"value"     db:"value"`
	ImageURL string  `json:"image_url" db:"image_url"`
	Set
}

type Set struct {
	Name   string `json:"set_name"   db:"set_name"`
	Series string `json:"set_series" db:"set_series"`
	Total  int
}
