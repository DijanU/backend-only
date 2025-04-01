package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Series struct {
	ID         int    `json:"id"`
	Ranking    int    `json:"ranking"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	Lwepisodes int    `json:"lwespisodes"`
	Tepisodes  int    `json:"tepisodes"`
}

func getallseries(w http.ResponseWriter, r *http.Request) {
	series := []Series{{Ranking: 1, Title: "Transformers", Status: "a", Lwepisodes: 7, Tepisodes: 14}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

func getseriesbyid() {}

func updateseiers() {}

func createSeries(w http.ResponseWriter, r *http.Request) {
	var series Series
	if err := json.NewDecoder(r.Body).Decode(&series); err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
	}
}

func deleteseries() {}

func main() {
	fmt.Println("blahhh")
	http.HandleFunc("/obtener", getallseries)
	http.HandleFunc("/crear", createSeries)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error al iniciar el s}ervidor:", err)
	}
}
