package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kevalsabhani/go-assignment/db"
	"github.com/kevalsabhani/go-assignment/handlers"
	"github.com/kevalsabhani/go-assignment/services"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	godotenv.Load()
	serverPort := os.Getenv("SERVER_PORT")

	// Connect DB
	db := db.ConnectDB()

	// Setup routes
	r := mux.NewRouter()

	// initialize handlers
	employeeHandler := handlers.NewEmployeeHandler(services.NewEmployeeService(db))
	handlers.SetupRoutes(r, employeeHandler)

	log.Printf("Server is running on :%s", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r))
}
