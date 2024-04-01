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

type CharSkillsController struct {
	db *sql.DB
}

// IntArrayScanner represents a custom type to scan array elements into integers
type IntArrayScanner []int

// Scan implements the sql.Scanner interface for IntArrayScanner
func (a *IntArrayScanner) Scan(src interface{}) error {
	if src == nil {
		*a = nil // Set to nil slice if the source is NULL
		return nil
	}

	str := string(src.([]byte))
	if str == "{}" {
		*a = []int{} // Set to empty slice if the source is an empty array
		return nil
	}

	parts := strings.Split(str[1:len(str)-1], ",")
	var result []int
	for _, part := range parts {
		i, err := strconv.Atoi(part)
		if err != nil {
			return err
		}
		result = append(result, i)
	}
	*a = result
	return nil
}

func NewCharSkillsController(db *sql.DB) *CharSkillsController {
	return &CharSkillsController{
		db: db,
	}
}

func (cc *CharSkillsController) GetAll(w http.ResponseWriter, r *http.Request) {
	lists, err := cc.getAllLists()
	if err != nil {
		log.Printf("Error querying all lists: %s", err)
		http.Error(w, fmt.Sprintf("Error getting lists: %s", err), http.StatusInternalServerError)
		return
	}

	// Convert lists to JSON and send it in the response
	responseJSON, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding lists to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// Implement this method to retrieve all lists from the database
func (cc *CharSkillsController) getAllLists() ([]CharSkill, error) {
	// Implement the logic to fetch all lists from the database
	// Example:
	rows, err := cc.db.Query("SELECT * FROM character_skills")
	if err != nil {
		log.Printf("Error querying all lists: %s", err)
		return nil, err
	}
	defer rows.Close()

	var lists []CharSkill

	for rows.Next() {
		var list CharSkill

		err := rows.Scan(&list.ID, &list.Name, &list.CharID, (*IntArrayScanner)(&list.SpellList), (*IntArrayScanner)(&list.CAList),
			(*IntArrayScanner)(&list.Boons), (*IntArrayScanner)(&list.Banes), &list.Budding, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lists = append(lists, list)
	}

	return lists, nil
}
func (cc *CharSkillsController) GetOneByID(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["listID"]
	id, err := strconv.Atoi(listID)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	character, err := cc.getListByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting list: %s", err), http.StatusInternalServerError)
		return
	}

	if character == nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(character)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding character to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CharSkillsController) getListByID(id int) (*CharSkill, error) {
	// Implement the logic to fetch a character by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM character_skills WHERE id = $1", id)

	var list CharSkill
	err := row.Scan(&list.ID, &list.Name, &list.CharID, (*IntArrayScanner)(&list.SpellList), (*IntArrayScanner)(&list.CAList),
		(*IntArrayScanner)(&list.Boons), (*IntArrayScanner)(&list.Banes), &list.Budding, &list.CreatedAt, &list.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &list, nil
}
func (cc *CharSkillsController) GetOneByCharID(w http.ResponseWriter, r *http.Request) {
	charID := mux.Vars(r)["charID"]
	id, err := strconv.Atoi(charID)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	character, err := cc.getListByCharID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting character: %s", err), http.StatusInternalServerError)
		return
	}

	if character == nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(character)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding character to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CharSkillsController) getListByCharID(id int) (*CharSkill, error) {
	// Implement the logic to fetch a character by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM character_skills WHERE char_id = $1", id)

	var list CharSkill
	err := row.Scan(&list.ID, &list.Name, &list.CharID, (*IntArrayScanner)(&list.SpellList), (*IntArrayScanner)(&list.CAList),
		(*IntArrayScanner)(&list.Boons), (*IntArrayScanner)(&list.Banes), &list.Budding, &list.CreatedAt, &list.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &list, nil
}

func (cc *CharSkillsController) PostOne(w http.ResponseWriter, r *http.Request) {
	var list CharSkill
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()

	err = cc.insertCharSkillList(&list)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting character skill list: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(list)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding list to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CharSkillsController) insertCharSkillList(list *CharSkill) error {
	// Perform the insert operation with the RETURNING clause to get the ID

	err := cc.db.QueryRow(`
		INSERT INTO character_skills (name, char_id, spell_list, ca_list, boons, banes, budding_talent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, list.Name, list.CharID, pq.Array(list.SpellList), pq.Array(list.CAList), pq.Array(list.Boons),
		pq.Array(list.Banes), list.Budding, list.CreatedAt,
		list.UpdatedAt).Scan(&list.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *CharSkillsController) PutOne(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["listID"]
	id, err := strconv.Atoi(listID)
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	var updatedList CharSkill
	err = json.NewDecoder(r.Body).Decode(&updatedList)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	updatedList.UpdatedAt = time.Now()

	err = cc.updateList(id, &updatedList)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating list: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(updatedList)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated list to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CharSkillsController) updateList(id int, updatedList *CharSkill) error {
	// Start building the SQL query
	query := "UPDATE character_skills SET updated_at = $1"
	args := []interface{}{updatedList.UpdatedAt}

	// Conditionally include fields in the update query
	if updatedList.Name != "" {
		query += ", name = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedList.Name)
	}
	if updatedList.CharID != 0 {
		query += ", char_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedList.CharID)
	}
	if len(updatedList.SpellList) != 0 {
		pqArray := pq.Array(updatedList.SpellList)
		query += ", spell_list = $" + strconv.Itoa(len(args)+1)
		args = append(args, pqArray)
	}
	if len(updatedList.CAList) != 0 {
		pqArray := pq.Array(updatedList.CAList)
		query += ", ca_list = $" + strconv.Itoa(len(args)+1)
		args = append(args, pqArray)
	}
	if updatedList.Boons != nil {
		pqArray := pq.Array(updatedList.Boons)
		query += ", boons = $" + strconv.Itoa(len(args)+1)
		args = append(args, pqArray)
	}
	if updatedList.Banes != nil {
		pqArray := pq.Array(updatedList.Banes)
		query += ", banes = $" + strconv.Itoa(len(args)+1)
		args = append(args, pqArray)
	}
	if updatedList.Budding != nil {
		query += ", budding_talent = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedList.Budding)
	}

	// Finish the query with the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)
	err := cc.db.QueryRow(query+" RETURNING *", args...).Scan(
		&updatedList.ID, &updatedList.Name, &updatedList.CharID,
		(*IntArrayScanner)(&updatedList.SpellList), (*IntArrayScanner)(&updatedList.CAList), (*IntArrayScanner)(&updatedList.Boons),
		(*IntArrayScanner)(&updatedList.Banes), &updatedList.Budding, &updatedList.CreatedAt,
		&updatedList.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cc *CharSkillsController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["listID"]
	id, err := strconv.Atoi(listID)
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	err = cc.deleteList(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("List with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting list: %s", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "msg": "Character deleted successfully."}`))
}

func (cc *CharSkillsController) deleteList(id int) error {
	// Check if the character exists
	var exists bool
	err := cc.db.QueryRow("SELECT EXISTS (SELECT 1 FROM character_skills WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("list with ID %d not found", id)
	}

	// Character exists, proceed with deletion
	_, err = cc.db.Exec("DELETE FROM character_skills WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
