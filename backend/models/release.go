package models

type Release struct {
	Id      int32    `json:"id"`
	Artists []string `json:"artists"`
	Styles  []string `json:"styles"`
	Genres  []string `json:"genres"`
}
