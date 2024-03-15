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

type SpellsController struct {
	db *sql.DB
}

func NewSpellsController(db *sql.DB) *SpellsController {
	return &SpellsController{
		db: db,
	}
}

func (cc *SpellsController) GetAll(w http.ResponseWriter, r *http.Request) {
	spells, err := cc.getAllSpells()
	if err != nil {
		log.Printf("Error querying all spells: %s", err)
		http.Error(w, fmt.Sprintf("Error getting spells: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(spells)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding spells to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SpellsController) getAllSpells() ([]Spells, error) {

	rows, err := cc.db.Query("SELECT * FROM spells")
	if err != nil {
		log.Printf("Error querying all spells: %s", err)
		return nil, err
	}
	defer rows.Close()

	var spells []Spells

	for rows.Next() {
		var spell Spells
		err := rows.Scan(&spell.ID, &spell.Name, &spell.Type, &spell.Might, &spell.Hit, &spell.Critical, &spell.Uses, &spell.Weight, &spell.RangeMin, &spell.RangeMax, &spell.Description, &spell.CreatedAt, &spell.UpdatedAt)
		if err != nil {
			return nil, err
		}
		spells = append(spells, spell)
	}

	return spells, nil
}

func (cc *SpellsController) GetOne(w http.ResponseWriter, r *http.Request) {
	spellID := mux.Vars(r)["spellID"]
	id, err := strconv.Atoi(spellID)
	if err != nil {
		http.Error(w, "Invalid spell ID", http.StatusBadRequest)
		return
	}

	spell, err := cc.getSpellByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting spell: %s", err), http.StatusInternalServerError)
		return
	}

	if spell == nil {
		http.Error(w, "Spell not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(spell)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding spell to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SpellsController) getSpellByID(id int) (*Spells, error) {
	// Implement the logic to fetch a spell by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM spells WHERE id = $1", id)

	var spell Spells
	err := row.Scan(&spell.ID, &spell.Name, &spell.Type, &spell.Might, &spell.Hit, &spell.Critical, &spell.Uses, &spell.Weight, &spell.RangeMin, &spell.RangeMax, &spell.Description, &spell.CreatedAt, &spell.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &spell, nil
}

func (cc *SpellsController) PostOne(w http.ResponseWriter, r *http.Request) {
	var spell Spells
	err := json.NewDecoder(r.Body).Decode(&spell)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	spell.CreatedAt = time.Now()
	spell.UpdatedAt = time.Now()

	err = cc.insertSpell(&spell)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting spell: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(spell)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding spell to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SpellsController) insertSpell(spell *Spells) error {
	// Perform the insert operation with the RETURNING clause to get the ID
	err := cc.db.QueryRow(`
		INSERT INTO spells (name, type, might, hit, critical, uses, weight, range_min, range_max, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`, spell.Name, spell.Type, spell.Might, spell.Hit, spell.Critical, spell.Uses, spell.Weight, spell.RangeMin, spell.RangeMax, spell.Description, spell.CreatedAt, spell.UpdatedAt).Scan(&spell.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *SpellsController) PutOne(w http.ResponseWriter, r *http.Request) {
	spellID := mux.Vars(r)["spellID"]
	id, err := strconv.Atoi(spellID)
	if err != nil {
		http.Error(w, "Invalid spell ID", http.StatusBadRequest)
		return
	}

	var updatedSpell Spells
	err = json.NewDecoder(r.Body).Decode(&updatedSpell)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	updatedSpell.UpdatedAt = time.Now()

	err = cc.updateSpell(id, &updatedSpell)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating spell: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(updatedSpell)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated spell to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SpellsController) updateSpell(id int, updatedSpell *Spells) error {
	// Start building the SQL query
	query := "UPDATE spells SET updated_at = $1"
	args := []interface{}{updatedSpell.UpdatedAt}

	// Conditionally include fields in the update query
	if updatedSpell.Name != "" {
		query += ", name = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Name)
	}
	if updatedSpell.Type != "" {
		query += ", type = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Type)
	}
	if updatedSpell.Might != nil {
		query += ", might = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Might)
	}
	if updatedSpell.Hit != nil {
		query += ", hit = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Hit)
	}
	if updatedSpell.Critical != nil {
		query += ", critical = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Critical)
	}
	if updatedSpell.Uses != 0 {
		query += ", uses = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Uses)
	}
	if updatedSpell.Weight != nil {
		query += ", weight = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Weight)
	}
	if updatedSpell.RangeMin != 0 {
		query += ", range_min = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.RangeMin)
	}
	if updatedSpell.RangeMax != nil {
		query += ", range_max = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.RangeMax)
	}
	if updatedSpell.Description != nil {
		query += ", description = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSpell.Description)
	}

	// Finish the query with the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)
	err := cc.db.QueryRow(query+" RETURNING *", args...).Scan(
		&updatedSpell.ID, &updatedSpell.Name, &updatedSpell.Type,
		&updatedSpell.Might, &updatedSpell.Hit, &updatedSpell.Critical,
		&updatedSpell.Uses, &updatedSpell.Weight, &updatedSpell.RangeMin,
		&updatedSpell.RangeMax, &updatedSpell.Description, &updatedSpell.CreatedAt,
		&updatedSpell.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cc *SpellsController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	spellID := mux.Vars(r)["spellID"]
	id, err := strconv.Atoi(spellID)
	if err != nil {
		http.Error(w, "Invalid spell ID", http.StatusBadRequest)
		return
	}

	err = cc.deleteSpell(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("Spell with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting spell: %s", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "msg": "Spell deleted successfully."}`))
}

func (cc *SpellsController) deleteSpell(id int) error {
	// Check if the spell exists
	var exists bool
	err := cc.db.QueryRow("SELECT EXISTS (SELECT 1 FROM spells WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("spell with ID %d not found", id)
	}

	// Spell exists, proceed with deletion
	_, err = cc.db.Exec("DELETE FROM spells WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
