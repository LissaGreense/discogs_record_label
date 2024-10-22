package models

type UniqueName struct {
	Name string `json:"name"`
}

type NameCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type CountResult struct {
	ReleaseCount int         `json:"releaseCount"`
	ArtistCounts []NameCount `json:"artistCounts"`
	StyleCounts  []NameCount `json:"styleCounts"`
	GenreCounts  []NameCount `json:"genreCounts"`
}
