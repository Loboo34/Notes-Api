package main

import (
	"fmt"
	"log"
	"net/http"

	"notes/database"
	"notes/handlers"
	"notes/logger"
)

func main() {

	logger.Logger()
	defer logger.Log.Sync()

	db := database.ConnectDB()
	fmt.Println("Using DB:", db.Name())

	mux := http.NewServeMux()

	mux.HandleFunc("/add", handlers.AddNote)
	mux.HandleFunc("/update", handlers.UpdateNote)
	mux.HandleFunc("/delete", handlers.DeleteNote)
	mux.HandleFunc("/get", handlers.GetNotes)
	mux.HandleFunc("/gett", handlers.GetNoteById)

	// Serve frontend files
	handler := enableCORS(mux)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
