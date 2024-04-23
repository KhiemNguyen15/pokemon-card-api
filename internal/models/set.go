package models

type Set struct {
	Name   string `json:"name"   db:"name"`
	Series string `json:"series" db:"series"`
	Total  int    `json:"total"  db:"total"`
}
