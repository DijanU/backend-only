package models

type Series struct {
	ID          int    `json:"id"`
	Ranking     int    `json:"ranking"`
	Title       string `json:"title"`
	Status      string `json:"status,omitempty"`
	LwsEpisodes int    `json:"lastEpisodeWatched"` // Cambiado a LwsEpisodes
	TEpisodes   int    `json:"totalEpisodes"`
}
