package main

import (
	"database/sql"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Retrieve environment variables
	host := os.Getenv("HOST")
	sql_port := os.Getenv("SQL_PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, sql_port, user, password, dbname)
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
	combatArtController := NewCombatArtController(db)
	weaponsController := NewWeaponsController(db)

	// Define your routes
	r.HandleFunc("/characters", characterController.GetAll).Methods("GET")
	r.HandleFunc("/characters/{charID}", characterController.GetOne).Methods("GET")
	r.HandleFunc("/characters/house/{affinity}", characterController.GetByAffinity).Methods("GET")
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

	r.HandleFunc("/combat_arts", combatArtController.GetAll).Methods("GET")
	r.HandleFunc("/combat_arts/{artID}", combatArtController.GetOne).Methods("GET")
	r.HandleFunc("/combat_arts", combatArtController.PostOne).Methods("POST")
	r.HandleFunc("/combat_arts/{artID}", combatArtController.PutOne).Methods("PUT")
	r.HandleFunc("/combat_arts/{artID}", combatArtController.DeleteOne).Methods("DELETE")

	r.HandleFunc("/weapons", weaponsController.GetAll).Methods("GET")
	r.HandleFunc("/weapons/{weaponID}", weaponsController.GetOne).Methods("GET")
	r.HandleFunc("/weapons/name/{weaponName}", weaponsController.GetOneName).Methods("GET")
	r.HandleFunc("/weapons", weaponsController.PostOne).Methods("POST")
	r.HandleFunc("/weapons/{weaponID}", weaponsController.PutOne).Methods("PUT")
	r.HandleFunc("/weapons/{weaponID}", weaponsController.DeleteOne).Methods("DELETE")

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "2999"
	}

	log.Printf("Server is running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
