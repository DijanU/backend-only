package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./series.db")
	if err != nil {
		return err
	}

	// Verificar que la tabla existe
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS series (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        ranking INTEGER UNIQUE NOT NULL,
        title TEXT NOT NULL,
        status TEXT,
        lws_episodes INTEGER DEFAULT 0,
        t_episodes INTEGER DEFAULT 0
    )`)

	return err
}

type Series struct {
	ID          int    `json:"id"`
	Ranking     int    `json:"ranking"`
	Title       string `json:"title"`
	Status      string `json:"status,omitempty"`
	LwsEpisodes int    `json:"lwespisodes"` // Cambiado a LwsEpisodes
	TEpisodes   int    `json:"tepisodes"`
}

func getallseries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM series")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var series []Series
	for rows.Next() {
		var s Series
		if err := rows.Scan(&s.ID, &s.Ranking, &s.Title, &s.Status, &s.LwsEpisodes, &s.TEpisodes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		series = append(series, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

func getseriesbyid() {}

func updateseiers() {}

func createSeries(w http.ResponseWriter, r *http.Request) {
	var newSeries Series
	if err := json.NewDecoder(r.Body).Decode(&newSeries); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validación básica
	if newSeries.Title == "" || newSeries.Ranking == 0 {
		http.Error(w, "Title and Ranking are required", http.StatusBadRequest)
		return
	}

	// Insertar en la base de datos
	result, err := db.Exec(
		"INSERT INTO series (ranking, title, status, lws_episodes, t_episodes) VALUES (?, ?, ?, ?, ?)",
		newSeries.Ranking,
		newSeries.Title,
		newSeries.Status,
		newSeries.LwsEpisodes, // Asegúrate de usar el mismo nombre
		newSeries.TEpisodes,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener el ID generado
	id, _ := result.LastInsertId()
	newSeries.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSeries)
}

func deleteseries() {}

func main() {
	// Inicializar la base de datos
	if err := initDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	router := chi.NewRouter()

	// Middleware CORS (versión mejorada)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Logging para desarrollo
			fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}) // <--- ¡Aquí estaba faltando el paréntesis!

	// Rutas API
	router.Route("/api", func(r chi.Router) {
		r.Get("/series", getallseries)
		r.Post("/series", createSeries)
		//router.Get("/series/:id", getseriesbyid)
		//router.Put("/series/:id", updateseiers)
		//router.Delete("/series/:id", deleteseries)

		//router.Patch("series/:id/upvote")
		//router.Patch("series/:id/downvote")
		//router.Patch("/series/:id/episode")
	})

	// Mensaje de inicio
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Available routes:")
	fmt.Println("  GET  /api/series")
	fmt.Println("  POST /api/series")

	// Iniciar servidor
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server error:", err)
	}
}
