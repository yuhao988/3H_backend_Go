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
	"github.com/lib/pq"
)

type ClassController struct {
	db *sql.DB
}

func NewClassController(db *sql.DB) *ClassController {
	return &ClassController{
		db: db,
	}
}

func (cc *ClassController) GetAll(w http.ResponseWriter, r *http.Request) {
	classes, err := cc.getAllClass()
	if err != nil {
		log.Printf("Error querying all classes: %s", err)
		http.Error(w, fmt.Sprintf("Error getting classes: %s", err), http.StatusInternalServerError)
		return
	}

	// Convert characters to JSON and send it in the response
	responseJSON, err := json.Marshal(classes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding classes to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// Implement this method to retrieve all characters from the database
func (cc *ClassController) getAllClass() ([]Classes, error) {
	// Implement the logic to fetch all characters from the database
	// Example:
	rows, err := cc.db.Query("SELECT * FROM classes")
	if err != nil {
		log.Printf("Error querying all classes: %s", err)
		return nil, err
	}
	defer rows.Close()

	var classes []Classes

	for rows.Next() {
		var class Classes
		err := rows.Scan(&class.ID, &class.Name, &class.Rank, (*IntArrayScanner)(&class.Base), (*IntArrayScanner)(&class.Bonus),
			(*IntArrayScanner)(&class.Growth), &class.CreatedAt, &class.UpdatedAt)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (cc *ClassController) GetOne(w http.ResponseWriter, r *http.Request) {
	classID := mux.Vars(r)["classID"]
	id, err := strconv.Atoi(classID)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	class, err := cc.getClassByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting class: %s", err), http.StatusInternalServerError)
		return
	}

	if class == nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(class)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding class to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *ClassController) getClassByID(id int) (*Classes, error) {
	// Implement the logic to fetch a character by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM classes WHERE id = $1", id)

	var class Classes
	err := row.Scan(&class.ID, &class.Name, &class.Rank, (*IntArrayScanner)(&class.Base),
		(*IntArrayScanner)(&class.Bonus), (*IntArrayScanner)(&class.Growth), &class.CreatedAt, &class.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &class, nil
}

func (cc *ClassController) PostOne(w http.ResponseWriter, r *http.Request) {
	var class Classes
	err := json.NewDecoder(r.Body).Decode(&class)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	class.CreatedAt = time.Now()
	class.UpdatedAt = time.Now()

	err = cc.insertClass(&class)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting class: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(class)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding class to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *ClassController) insertClass(class *Classes) error {
	// Perform the insert operation with the RETURNING clause to get the ID

	err := cc.db.QueryRow(`
		INSERT INTO classes (name, rank, base, bonus, growth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, class.Name, class.Rank, pq.Array(class.Base), pq.Array(class.Bonus), pq.Array(class.Growth),
		class.CreatedAt, class.UpdatedAt).Scan(&class.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *ClassController) PutOne(w http.ResponseWriter, r *http.Request) {
	classID := mux.Vars(r)["classID"]
	id, err := strconv.Atoi(classID)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	var updatedClass Classes
	err = json.NewDecoder(r.Body).Decode(&updatedClass)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	updatedClass.UpdatedAt = time.Now()

	err = cc.updateClass(id, &updatedClass)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating class: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(updatedClass)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated class to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *ClassController) updateClass(id int, updatedClass *Classes) error {
	// Start building the SQL query
	query := "UPDATE classes SET updated_at = $1"
	args := []interface{}{updatedClass.UpdatedAt}

	// Conditionally include fields in the update query
	if updatedClass.Name != "" {
		query += ", name = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedClass.Name)
	}
	if updatedClass.Rank != "" {
		query += ", rank = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedClass.Rank)
	}
	if len(updatedClass.Base) != 0 {
		pqArray := pq.Array(updatedClass.Base)
		query += ", base = $" + strconv.Itoa(len(args)+1)
		args = append(args, pqArray)
	}
	if len(updatedClass.Bonus) != 0 {
		pqArray := pq.Array(updatedClass.Bonus)
		query += ", bonus = $" + strconv.Itoa(len(args)+1)
		args = append(args, pqArray)
	}
	if updatedClass.Growth != nil {
		pqArray := pq.Array(updatedClass.Growth)
		query += ", growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, pqArray)
	}

	// Finish the query with the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)
	err := cc.db.QueryRow(query+" RETURNING *", args...).Scan(
		&updatedClass.ID, &updatedClass.Name, &updatedClass.Rank,
		(*IntArrayScanner)(&updatedClass.Base), (*IntArrayScanner)(&updatedClass.Bonus),
		(*IntArrayScanner)(&updatedClass.Growth), &updatedClass.CreatedAt,
		&updatedClass.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cc *ClassController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	charID := mux.Vars(r)["charID"]
	id, err := strconv.Atoi(charID)
	if err != nil {
		http.Error(w, "invalid character ID", http.StatusBadRequest)
		return
	}

	err = cc.deleteClass(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("Class with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting character: %s", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "msg": "Class deleted successfully."}`))
}

func (cc *ClassController) deleteClass(id int) error {
	// Check if the character exists
	var exists bool
	err := cc.db.QueryRow("SELECT EXISTS (SELECT 1 FROM characters WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("character with ID %d not found", id)
	}

	// Class exists, proceed with deletion
	_, err = cc.db.Exec("DELETE FROM characters WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
