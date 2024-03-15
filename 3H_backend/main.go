package main

import (
	"database/sql"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "921022"
	dbname   = "fe3h"
)

func main() {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")

	// Set up the Gorilla mux router
	r := mux.NewRouter()

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	characterController := NewCharacterController(db)
	skillsController := NewSkillsController(db)
	spellsController := NewSpellsController(db)

	// Define your routes
	r.HandleFunc("/characters", characterController.GetAll).Methods("GET")
	r.HandleFunc("/characters/{charID}", characterController.GetOne).Methods("GET")
	r.HandleFunc("/characters", characterController.PostOne).Methods("POST")
	r.HandleFunc("/characters/{charID}", characterController.PutOne).Methods("PUT")
	r.HandleFunc("/characters/{charID}", characterController.DeleteOne).Methods("DELETE")

	r.HandleFunc("/skill_types", skillsController.GetAll).Methods("GET")
	r.HandleFunc("/skill_types/{skillID}", skillsController.GetOne).Methods("GET")
	r.HandleFunc("/skill_types", skillsController.PostOne).Methods("POST")
	r.HandleFunc("/skill_types/{skillID}", skillsController.PutOne).Methods("PUT")
	r.HandleFunc("/skill_types/{skillID}", skillsController.DeleteOne).Methods("DELETE")

	r.HandleFunc("/spells", spellsController.GetAll).Methods("GET")
	r.HandleFunc("/spells/{spellID}", spellsController.GetOne).Methods("GET")
	r.HandleFunc("/spells", spellsController.PostOne).Methods("POST")
	r.HandleFunc("/spells/{spellID}", spellsController.PutOne).Methods("PUT")
	r.HandleFunc("/spells/{spellID}", spellsController.DeleteOne).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "2999"
	}

	log.Printf("Server is running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
