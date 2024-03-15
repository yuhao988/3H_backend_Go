package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type CombatArtController struct {
	db *sql.DB
}

func NewCombatArtController(db *sql.DB) *CombatArtController {
	return &CombatArtController{
		db: db,
	}
}

func (cc *CombatArtController) GetAll(w http.ResponseWriter, r *http.Request) {
	combatArts, err := cc.getAllCombatArts()
	if err != nil {
		log.Printf("Error querying all combat arts: %s", err)
		http.Error(w, fmt.Sprintf("Error getting combat arts: %s", err), http.StatusInternalServerError)
		return
	}

	// Convert combatArts to JSON and send it in the response
	responseJSON, err := json.Marshal(combatArts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding combat arts to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// Implement this method to retrieve all combatArts from the database
func (cc *CombatArtController) getAllCombatArts() ([]CombatArts, error) {
	// Implement the logic to fetch all combatArts from the database
	// Example:
	rows, err := cc.db.Query("SELECT * FROM combat_arts")
	if err != nil {
		log.Printf("Error querying all combat arts: %s", err)
		return nil, err
	}
	defer rows.Close()

	var combatArts []CombatArts

	for rows.Next() {
		var combatArt CombatArts
		err := rows.Scan(&combatArt.ID, &combatArt.Name, &combatArt.TypeID, &combatArt.StrMag,
			&combatArt.Might, &combatArt.Hit, &combatArt.Critical, &combatArt.DurabilityCost,
			&combatArt.RangeMin, &combatArt.RangeMax, &combatArt.CreatedAt, &combatArt.UpdatedAt)
		if err != nil {
			return nil, err
		}
		combatArts = append(combatArts, combatArt)
	}

	return combatArts, nil
}

func (cc *CombatArtController) GetOne(w http.ResponseWriter, r *http.Request) {
	artID := mux.Vars(r)["artID"]
	id, err := strconv.Atoi(artID)
	if err != nil {
		http.Error(w, "Invalid combat art ID", http.StatusBadRequest)
		return
	}

	combatArt, err := cc.getCombatArtByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting combat art: %s", err), http.StatusInternalServerError)
		return
	}

	if combatArt == nil {
		http.Error(w, "combat art not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(combatArt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding combat art to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CombatArtController) getCombatArtByID(id int) (*CombatArts, error) {
	// Implement the logic to fetch a combatArt by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM combat_arts WHERE id = $1", id)

	var combatArt CombatArts
	err := row.Scan(&combatArt.ID, &combatArt.Name, &combatArt.TypeID, &combatArt.Might, &combatArt.Hit,
		&combatArt.Critical, &combatArt.DurabilityCost, &combatArt.RangeMin, &combatArt.RangeMax,
		&combatArt.Description, &combatArt.CreatedAt, &combatArt.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &combatArt, nil
}

func (cc *CombatArtController) PostOne(w http.ResponseWriter, r *http.Request) {
	var combatArt CombatArts
	err := json.NewDecoder(r.Body).Decode(&combatArt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	combatArt.CreatedAt = time.Now()
	combatArt.UpdatedAt = time.Now()

	err = cc.insertCombatArt(&combatArt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting combat art: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(combatArt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding combat art to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CombatArtController) insertCombatArt(combatArt *CombatArts) error {
	// Perform the insert operation with the RETURNING clause to get the ID
	err := cc.db.QueryRow(`
		INSERT INTO combat_arts (name, type_id, str_mag, might, hit, critical, durability_cost,
			 range_min, range_max, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,$12)
		RETURNING id
	`, combatArt.Name, combatArt.TypeID, combatArt.StrMag, combatArt.Might, combatArt.Hit, combatArt.Critical, combatArt.DurabilityCost,
		combatArt.RangeMin, combatArt.RangeMax, combatArt.Description, combatArt.CreatedAt, combatArt.UpdatedAt).Scan(&combatArt.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *CombatArtController) PutOne(w http.ResponseWriter, r *http.Request) {
	artID := mux.Vars(r)["artID"]
	id, err := strconv.Atoi(artID)
	if err != nil {
		http.Error(w, "Invalid combat art ID", http.StatusBadRequest)
		return
	}

	var updatedArt CombatArts
	err = json.NewDecoder(r.Body).Decode(&updatedArt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	updatedArt.UpdatedAt = time.Now()

	err = cc.updatedArt(id, &updatedArt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating combat art: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(updatedArt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated combat art to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CombatArtController) updatedArt(id int, updatedArt *CombatArts) error {
	// Start building the SQL query
	query := "UPDATE combat_arts SET updated_at = $1"
	args := []interface{}{updatedArt.UpdatedAt}

	// Conditionally include fields in the update query
	if updatedArt.Name != "" {
		query += ", name = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.Name)
	}
	if updatedArt.TypeID != 0 {
		query += ", type_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.TypeID)
	}
	if updatedArt.StrMag != nil {
		query += ", str_mag = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.StrMag)
	}
	if updatedArt.Might != nil {
		query += ", might = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.Might)
	}
	if updatedArt.Hit != nil {
		query += ", hit = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.Hit)
	}
	if updatedArt.Critical != nil {
		query += ", critical = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.Critical)
	}
	if updatedArt.DurabilityCost != 0 {
		query += ", durability_cost = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.DurabilityCost)
	}

	if updatedArt.RangeMin != 0 {
		query += ", range_min = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.RangeMin)
	}
	if updatedArt.RangeMax != nil {
		query += ", range_max = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.RangeMax)
	}
	if updatedArt.Description != nil {
		query += ", description = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedArt.Description)
	}

	// Finish the query with the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)
	err := cc.db.QueryRow(query+" RETURNING *", args...).Scan(
		&updatedArt.ID, &updatedArt.Name, &updatedArt.TypeID, &updatedArt.StrMag,
		&updatedArt.Might, &updatedArt.Hit, &updatedArt.Critical,
		&updatedArt.DurabilityCost, &updatedArt.RangeMin,
		&updatedArt.RangeMax, &updatedArt.Description, &updatedArt.CreatedAt,
		&updatedArt.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cc *CombatArtController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	artID := mux.Vars(r)["artID"]
	id, err := strconv.Atoi(artID)
	if err != nil {
		http.Error(w, "Invalid combat art ID", http.StatusBadRequest)
		return
	}

	err = cc.deleteComabtArt(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("Combat art with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting combat art: %s", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "msg": "Combat art deleted successfully."}`))
}

func (cc *CombatArtController) deleteComabtArt(id int) error {
	// Check if the combat art exists
	var exists bool
	err := cc.db.QueryRow("SELECT EXISTS (SELECT 1 FROM combat_arts WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Combat art with ID %d not found", id)
	}

	// Combat art exists, proceed with deletion
	_, err = cc.db.Exec("DELETE FROM combat_arts WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
