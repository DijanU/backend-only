// @title Series API
// @version 1.0.0
// @description This is a sample API for managing TV series.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/series
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/DijanU/backend-only/docs"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
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

// estructura de cada serie
type Series struct {
	ID          int    `json:"id"`
	Ranking     int    `json:"ranking"`
	Title       string `json:"title"`
	Status      string `json:"status,omitempty"`
	LwsEpisodes int    `json:"lastEpisodeWatched"`
	TEpisodes   int    `json:"totalEpisodes"`
}

// @Summary Obtener todas las series
// @Description Retorna una lista de todas las series guardadas en la base de datos
// @Tags series
// @Produce json
// @Success 200 {array} Series
// @Router /series [get]
func getallseries(w http.ResponseWriter, r *http.Request) {
	//query select
	rows, err := db.Query("SELECT * FROM series")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	//array de series para exportar por el API
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

// @Summary Obtener una serie por ID
// @Description Retorna una serie específica por su ID
// @Tags series
// @Produce json
// @Param id path int true "ID de la serie"
// @Success 200 {object} Series
// @Failure 400 {string} string "ID inválido"
// @Failure 404 {string} string "Serie no encontrada"
// @Router /series/{id} [get]
func getseriesbyid(w http.ResponseWriter, r *http.Request) {

	//conertir id a un it desde el request
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var s Series

	//encontrar el regsitro de la serie conese id
	err = db.QueryRow("SELECT id, ranking, title, status, lws_episodes, t_episodes FROM series WHERE id = ?", id).
		Scan(&s.ID, &s.Ranking, &s.Title, &s.Status, &s.LwsEpisodes, &s.TEpisodes)

	if err == sql.ErrNoRows {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)

}

// @Summary Actualizar una serie
// @Description Actualiza la información de una serie existente
// @Tags series
// @Accept json
// @Produce json
// @Param id path int true "ID de la serie"
// @Param serie body Series true "Datos actualizados"
// @Success 200 {object} Series
// @Failure 400 {string} string "ID o JSON inválido"
// @Failure 404 {string} string "Serie no encontrada"
// @Router /series/{id} [put]
func updateseiers(w http.ResponseWriter, r *http.Request) {

	//conertir id a un it desde el request
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var s Series
	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	s.ID = id
	//actualizar registro de bd
	updated, err := db.Exec("UPDATE series SET ranking = ?, title = ?, status = ?, lws_episodes = ?, t_episodes = ?  WHERE id = ?", s.Ranking, s.Title, s.Status, s.LwsEpisodes, s.TEpisodes, s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar cuántas filas fueron afectadas
	rowsAffected, _ := updated.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)

}

// @Summary Crear una nueva serie
// @Description Crea una nueva serie con los datos enviados
// @Tags series
// @Accept json
// @Produce json
// @Param serie body Series true "Datos de la nueva serie"
// @Success 201 {object} Series
// @Failure 400 {string} string "Datos inválidos"
// @Router /series [post]
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

// @Summary Eliminar una serie
// @Description Elimina una serie existente por su ID
// @Tags series
// @Param id path int true "ID de la serie"
// @Success 204 {string} string "Sin contenido"
// @Failure 400 {string} string "ID inválido"
// @Failure 404 {string} string "Serie no encontrada"
// @Router /series/{id} [delete]
func deleteseries(w http.ResponseWriter, r *http.Request) {
	//conertir id a un it desde el request
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	//delete wehre id coincide
	deleted, err := db.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar cuántas filas fueron afectadas
	rowsAffected, _ := deleted.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Votar por una serie (upvote)
// @Description Incrementa el ranking de una serie específica
// @Tags series
// @Param id path int true "ID de la serie"
// @Success 200 {string} string "Voto aplicado"
// @Failure 400 {string} string "ID inválido"
// @Router /series/{id}/upvote [post]
func seriesupvote(w http.ResponseWriter, r *http.Request) {
	//conertir id a un it desde el request
	idStr := chi.URLParam(r, "id")

	// Convertir el ID de string a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Actualizar el ranking, sumando 1 al valor actual
	upvoted, err := db.Exec("UPDATE series SET ranking = ranking + 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar cuántas filas fueron afectadas
	rowsAffected, _ := upvoted.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Votar por una serie (downvote)
// @Description Disminuye el ranking de una serie específica
// @Tags series
// @Param id path int true "ID de la serie"
// @Success 200 {string} string "Voto aplicado"
// @Failure 400 {string} string "ID inválido"
// @Router /series/{id}/upvote [post]
func seriesdownvote(w http.ResponseWriter, r *http.Request) {
	//conertir id a un it desde el request
	idStr := chi.URLParam(r, "id")

	// Convertir el ID de string a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Actualizar el ranking, sumando 1 al valor actual
	upvoted, err := db.Exec("UPDATE series SET ranking = ranking - 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar cuántas filas fueron afectadas
	rowsAffected, _ := upvoted.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	// Responder con un código 204 (sin contenido)
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Sumar un episodio visto
// @Description Incrementa en 1 el número de episodios vistos (`lws_episodes`) de una serie específica
// @Tags series
// @Param id path int true "ID de la serie"
// @Success 204 {string} string "Episodio sumado"
// @Failure 400 {string} string "ID inválido"
// @Failure 404 {string} string "Serie no encontrada"
// @Failure 500 {string} string "Error del servidor"
// @Router /series/{id}/episodeplusone [patch]
func episodeplusone(w http.ResponseWriter, r *http.Request) {
	//conertir id a un it desde el request
	idStr := chi.URLParam(r, "id")

	// Convertir el ID de string a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Actualizar el ranking, sumando 1 al valor actual
	upvoted, err := db.Exec("UPDATE series SET lws_episodes = lws_episodes + 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar cuántas filas fueron afectadas
	rowsAffected, _ := upvoted.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Cambiar estado de una serie
// @Description Actualiza el campo `status` de una serie con el valor enviado en el cuerpo
// @Tags series
// @Param id path int true "ID de la serie"
// @Param status body object true "Nuevo estado en formato JSON. Ej: {\"status\": \"Completada\"}"
// @Success 204 {string} string "Estado actualizado"
// @Failure 400 {string} string "ID inválido o JSON mal formado"
// @Failure 404 {string} string "Serie no encontrada"
// @Failure 500 {string} string "Error del servidor"
// @Router /series/{id}/status [patch]
func statuschange(w http.ResponseWriter, r *http.Request) {
	//conertir id a un it desde el request
	idStr := chi.URLParam(r, "id")

	// Convertir el ID de string a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	//struct para recibir el nuevo estatus
	var status struct {
		Status string `json:"status"`
	}

	err = json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Actualizar el ranking, sumando 1 al valor actual
	upvoted, err := db.Exec("UPDATE series SET status = ? WHERE id = ?", status.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar cuántas filas fueron afectadas
	rowsAffected, _ := upvoted.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	// Responder con un código 204 (sin contenido)
	w.WriteHeader(http.StatusNoContent)
}

// @title API de Series
// @version 1.0
// @description Esta API permite manejar información de series, incluyendo episodios vistos, estatus, y votaciones.
// @host localhost:8080
// @BasePath /
func main() {
	// Inicializar la base de datos
	if err := initDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	router := chi.NewRouter()

	// Middleware CORS
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Logging para desarrollo
			fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Rutas API
	router.Route("/api/series", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Get("/", getallseries)
		r.Post("/", createSeries)
		r.Get("/{id}", getseriesbyid)
		r.Put("/{id}", updateseiers)
		r.Delete("/{id}", deleteseries)

		r.Patch("/{id}/upvote", seriesupvote)
		r.Patch("/{id}/downvote", seriesdownvote)
		r.Patch("/{id}/episode", episodeplusone)
		r.Patch("/{id}/status", statuschange)
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
