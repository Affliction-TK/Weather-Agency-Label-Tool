package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

// Models
type Station struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Image struct {
	ID         int       `json:"id"`
	Filename   string    `json:"filename"`
	Filepath   string    `json:"filepath"`
	UploadedAt time.Time `json:"uploaded_at"`
	Annotated  bool      `json:"annotated"`
}

type Annotation struct {
	ID              int       `json:"id"`
	ImageID         int       `json:"image_id"`
	Category        string    `json:"category"`
	Severity        string    `json:"severity"`
	ObservationTime time.Time `json:"observation_time"`
	Location        string    `json:"location"`
	Longitude       float64   `json:"longitude"`
	Latitude        float64   `json:"latitude"`
	StationID       int       `json:"station_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ImageWithAnnotation struct {
	Image      Image       `json:"image"`
	Annotation *Annotation `json:"annotation,omitempty"`
}

// Initialize database connection
func initDB() error {
	var err error
	// Connection string format: username:password@tcp(host:port)/dbname
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/weather_label_db?parseTime=true"
	}
	
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	
	if err = db.Ping(); err != nil {
		return err
	}
	
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	
	log.Println("Database connected successfully")
	return nil
}

// Calculate distance between two points using Haversine formula
func haversineDistance(lon1, lat1, lon2, lat2 float64) float64 {
	const earthRadius = 6371 // km
	
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	
	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

// Find nearest station to given coordinates
func findNearestStation(lon, lat float64) (*Station, error) {
	rows, err := db.Query("SELECT id, name, longitude, latitude FROM stations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var nearestStation *Station
	minDistance := math.MaxFloat64
	
	for rows.Next() {
		var station Station
		if err := rows.Scan(&station.ID, &station.Name, &station.Longitude, &station.Latitude); err != nil {
			continue
		}
		
		distance := haversineDistance(lon, lat, station.Longitude, station.Latitude)
		if distance < minDistance {
			minDistance = distance
			nearestStation = &station
		}
	}
	
	return nearestStation, nil
}

// API Handlers
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func getStations(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, longitude, latitude FROM stations ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	stations := []Station{}
	for rows.Next() {
		var station Station
		if err := rows.Scan(&station.ID, &station.Name, &station.Longitude, &station.Latitude); err != nil {
			continue
		}
		stations = append(stations, station)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)
}

func getNearestStation(w http.ResponseWriter, r *http.Request) {
	lonStr := r.URL.Query().Get("longitude")
	latStr := r.URL.Query().Get("latitude")
	
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}
	
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	
	station, err := findNearestStation(lon, lat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)
}

func getImages(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT i.id, i.filename, i.filepath, i.uploaded_at, i.annotated
		FROM images i
		ORDER BY i.annotated ASC, i.uploaded_at DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	images := []Image{}
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.ID, &img.Filename, &img.Filepath, &img.UploadedAt, &img.Annotated); err != nil {
			continue
		}
		images = append(images, img)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	var img Image
	err := db.QueryRow(`
		SELECT id, filename, filepath, uploaded_at, annotated
		FROM images
		WHERE id = ?
	`, id).Scan(&img.ID, &img.Filename, &img.Filepath, &img.UploadedAt, &img.Annotated)
	
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	
	// Get annotation if exists
	var annotation Annotation
	err = db.QueryRow(`
		SELECT id, image_id, category, severity, observation_time, location, 
		       longitude, latitude, station_id, created_at, updated_at
		FROM annotations
		WHERE image_id = ?
	`, img.ID).Scan(
		&annotation.ID, &annotation.ImageID, &annotation.Category, &annotation.Severity,
		&annotation.ObservationTime, &annotation.Location, &annotation.Longitude,
		&annotation.Latitude, &annotation.StationID, &annotation.CreatedAt, &annotation.UpdatedAt,
	)
	
	response := ImageWithAnnotation{
		Image: img,
	}
	
	if err == nil {
		response.Annotation = &annotation
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createAnnotation(w http.ResponseWriter, r *http.Request) {
	var annotation Annotation
	if err := json.NewDecoder(r.Body).Decode(&annotation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Check if annotation already exists for this image
	var existingID int
	err := db.QueryRow("SELECT id FROM annotations WHERE image_id = ?", annotation.ImageID).Scan(&existingID)
	
	if err == sql.ErrNoRows {
		// Create new annotation
		result, err := db.Exec(`
			INSERT INTO annotations (image_id, category, severity, observation_time, location, 
			                        longitude, latitude, station_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, annotation.ImageID, annotation.Category, annotation.Severity, annotation.ObservationTime,
			annotation.Location, annotation.Longitude, annotation.Latitude, annotation.StationID)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		id, _ := result.LastInsertId()
		annotation.ID = int(id)
	} else {
		// Update existing annotation
		_, err := db.Exec(`
			UPDATE annotations 
			SET category = ?, severity = ?, observation_time = ?, location = ?, 
			    longitude = ?, latitude = ?, station_id = ?
			WHERE image_id = ?
		`, annotation.Category, annotation.Severity, annotation.ObservationTime,
			annotation.Location, annotation.Longitude, annotation.Latitude, 
			annotation.StationID, annotation.ImageID)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		annotation.ID = existingID
	}
	
	// Mark image as annotated
	_, err = db.Exec("UPDATE images SET annotated = TRUE WHERE id = ?", annotation.ImageID)
	if err != nil {
		log.Printf("Error updating image annotated status: %v", err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(annotation)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := filepath.Join(uploadsDir, filename)
	
	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	
	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Save to database
	result, err := db.Exec("INSERT INTO images (filename, filepath) VALUES (?, ?)", filename, filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	id, _ := result.LastInsertId()
	
	img := Image{
		ID:         int(id),
		Filename:   filename,
		Filepath:   filepath,
		UploadedAt: time.Now(),
		Annotated:  false,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(img)
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	
	filepath := filepath.Join("./uploads", filename)
	
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	
	http.ServeFile(w, r, filepath)
}

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Create router
	r := mux.NewRouter()
	
	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/stations", getStations).Methods("GET")
	api.HandleFunc("/stations/nearest", getNearestStation).Methods("GET")
	api.HandleFunc("/images", getImages).Methods("GET")
	api.HandleFunc("/images/{id}", getImage).Methods("GET")
	api.HandleFunc("/annotations", createAnnotation).Methods("POST")
	api.HandleFunc("/upload", uploadImage).Methods("POST")
	
	// Image serving route
	r.HandleFunc("/images/{filename}", serveImage).Methods("GET")
	
	// Serve static files from frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist")))
	
	// Enable CORS
	handler := enableCORS(r)
	
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
