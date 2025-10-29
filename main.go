package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Loboo34/Notes-Api/database"
	"github.com/Loboo34/Notes-Api/handlers"
	"github.com/Loboo34/Notes-Api/logger"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// func main() {

// 	logger.Logger()
// 	defer logger.Log.Sync()

// 	db := database.ConnectDB()
// 	fmt.Println("Using DB:", db.Name())

// 	mux := http.NewServeMux()
// 	//user
// 	mux.HandleFunc("/register", handlers.RegisterUser)
// 	mux.HandleFunc("/login", handlers.LoginUser)

// 	//crud
// 	mux.HandleFunc("/add", handlers.AddNote)
// 	mux.HandleFunc("/update", handlers.UpdateNote)
// 	mux.HandleFunc("/delete", handlers.DeleteNote)
// 	mux.HandleFunc("/get", handlers.GetNotes)
// 	mux.HandleFunc("/gett", handlers.GetNoteById)

// 	// Serve frontend files
// 	handler := enableCORS(mux)

// 	fmt.Println("Server is running on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", handler))

// }

func main() {
	logger.Logger()
	defer logger.Log.Sync()
	db := database.ConnectDB()
	fmt.Println("Using Db:", db.Name())

	r := mux.NewRouter()

	r.Use(corsMiddleware())

	//auth
	r.HandleFunc("register", handlers.RegisterUser).Methods("Post")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	//crud
	r.HandleFunc("/add", handlers.AddNote).Methods("POST")
	r.HandleFunc("/update/{id}", handlers.UpdateNote).Methods("PUT")
	r.HandleFunc("delete/{id}", handlers.DeleteNote).Methods("DELETE")
	r.HandleFunc("/get/{id}", handlers.GetNoteById).Methods("GET")
	r.HandleFunc("/get", handlers.GetNotes).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+port, r))

}

func corsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
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
}
