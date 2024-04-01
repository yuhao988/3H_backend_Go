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

type CharacterController struct {
	db *sql.DB
}

func NewCharacterController(db *sql.DB) *CharacterController {
	return &CharacterController{
		db: db,
	}
}

func (cc *CharacterController) GetAll(w http.ResponseWriter, r *http.Request) {
	characters, err := cc.getAllCharacters()
	if err != nil {
		log.Printf("Error querying all characters: %s", err)
		http.Error(w, fmt.Sprintf("Error getting characters: %s", err), http.StatusInternalServerError)
		return
	}

	// Convert characters to JSON and send it in the response
	responseJSON, err := json.Marshal(characters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding characters to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// Implement this method to retrieve all characters from the database
func (cc *CharacterController) getAllCharacters() ([]Character, error) {
	// Implement the logic to fetch all characters from the database
	// Example:
	rows, err := cc.db.Query("SELECT * FROM characters")
	if err != nil {
		log.Printf("Error querying all characters: %s", err)
		return nil, err
	}
	defer rows.Close()

	var characters []Character

	for rows.Next() {
		var character Character
		err := rows.Scan(&character.ID, &character.Name, &character.ImageLink, &character.Affinity, &character.BaseLv, &character.HP, &character.HpGrowth, &character.Strength, &character.StrGrowth, &character.Magic, &character.MagGrowth, &character.Dexterity, &character.DexGrowth, &character.Speed, &character.SpdGrowth, &character.Luck, &character.LckGrowth, &character.Defence, &character.DefGrowth, &character.Resistance, &character.ResGrowth, &character.Charm, &character.ChaGrowth, &character.CreatedAt, &character.UpdatedAt)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}

func (cc *CharacterController) GetOne(w http.ResponseWriter, r *http.Request) {
	charID := mux.Vars(r)["charID"]
	id, err := strconv.Atoi(charID)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	character, err := cc.getCharacterByID(id)
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

func (cc *CharacterController) getCharacterByID(id int) (*Character, error) {
	// Implement the logic to fetch a character by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM characters WHERE id = $1", id)

	var character Character
	err := row.Scan(&character.ID, &character.Name, &character.ImageLink, &character.Affinity, &character.BaseLv, &character.HP, &character.HpGrowth, &character.Strength, &character.StrGrowth, &character.Magic, &character.MagGrowth, &character.Dexterity, &character.DexGrowth, &character.Speed, &character.SpdGrowth, &character.Luck, &character.LckGrowth, &character.Defence, &character.DefGrowth, &character.Resistance, &character.ResGrowth, &character.Charm, &character.ChaGrowth, &character.CreatedAt, &character.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &character, nil
}

func (cc *CharacterController) GetByAffinity(w http.ResponseWriter, r *http.Request) {
	affinity := mux.Vars(r)["affinity"]
	// aff, err := strconv.Atoi(affinity)
	// if err != nil {
	// 	http.Error(w, "Invalid character affinity", http.StatusBadRequest)
	// 	return
	// }

	character, err := cc.getCharacterByAffinity(affinity)
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

func (cc *CharacterController) getCharacterByAffinity(affinity string) ([]Character, error) {
	// Implement the logic to fetch a character by ID from the database
	// Example:
	rows, err := cc.db.Query("SELECT * FROM characters WHERE affinity = $1", affinity)
	if err != nil {
		log.Printf("Error querying all characters: %s", err)
		return nil, err
	}
	defer rows.Close()

	var characters []Character

	for rows.Next() {
		var character Character
		err := rows.Scan(&character.ID, &character.Name, &character.ImageLink, &character.Affinity, &character.BaseLv, &character.HP, &character.HpGrowth, &character.Strength, &character.StrGrowth, &character.Magic, &character.MagGrowth, &character.Dexterity, &character.DexGrowth, &character.Speed, &character.SpdGrowth, &character.Luck, &character.LckGrowth, &character.Defence, &character.DefGrowth, &character.Resistance, &character.ResGrowth, &character.Charm, &character.ChaGrowth, &character.CreatedAt, &character.UpdatedAt)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}
func (cc *CharacterController) GetByName(w http.ResponseWriter, r *http.Request) {
	charName := mux.Vars(r)["charName"]

	character, err := cc.getCharacterByName(charName)
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

func (cc *CharacterController) getCharacterByName(name string) (*Character, error) {
	name = strings.Title(name)
	// Implement the logic to fetch a character by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM characters WHERE name = $1", name)

	var character Character
	err := row.Scan(&character.ID, &character.Name, &character.ImageLink, &character.Affinity, &character.BaseLv, &character.HP, &character.HpGrowth, &character.Strength, &character.StrGrowth, &character.Magic, &character.MagGrowth, &character.Dexterity, &character.DexGrowth, &character.Speed, &character.SpdGrowth, &character.Luck, &character.LckGrowth, &character.Defence, &character.DefGrowth, &character.Resistance, &character.ResGrowth, &character.Charm, &character.ChaGrowth, &character.CreatedAt, &character.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &character, nil
}

func (cc *CharacterController) PostOne(w http.ResponseWriter, r *http.Request) {
	var character Character
	err := json.NewDecoder(r.Body).Decode(&character)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	character.CreatedAt = time.Now()
	character.UpdatedAt = time.Now()

	err = cc.insertCharacter(&character)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting character: %s", err), http.StatusInternalServerError)
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

func (cc *CharacterController) insertCharacter(character *Character) error {
	// Perform the insert operation with the RETURNING clause to get the ID
	err := cc.db.QueryRow(`
		INSERT INTO characters (name, image_link, affinity, base_lv, hp, hp_growth,
			 strength, str_growth, magic, mag_growth, dexterity, dex_growth, speed, 
			 spd_growth, luck, lck_growth, defence, def_growth, resistance, res_growth, 
			 charm, cha_growth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
			$17, $18, $19, $20, $21, $22, $23, $24)
		RETURNING id
	`, character.Name, character.ImageLink, character.Affinity, character.BaseLv,
		character.HP, character.HpGrowth, character.Strength, character.StrGrowth,
		character.Magic, character.MagGrowth, character.Dexterity, character.DexGrowth,
		character.Speed, character.SpdGrowth, character.Luck, character.LckGrowth,
		character.Defence, character.DefGrowth, character.Resistance, character.ResGrowth,
		character.Charm, character.ChaGrowth, character.CreatedAt, character.UpdatedAt).Scan(&character.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *CharacterController) PutOne(w http.ResponseWriter, r *http.Request) {
	charID := mux.Vars(r)["charID"]
	id, err := strconv.Atoi(charID)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	var updatedCharacter Character
	err = json.NewDecoder(r.Body).Decode(&updatedCharacter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	updatedCharacter.UpdatedAt = time.Now()

	err = cc.updateCharacter(id, &updatedCharacter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating character: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(updatedCharacter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated character to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *CharacterController) updateCharacter(id int, updatedCharacter *Character) error {
	// Start building the SQL query
	query := "UPDATE characters SET updated_at = $1"
	args := []interface{}{updatedCharacter.UpdatedAt}

	// Conditionally include fields in the update query
	if updatedCharacter.Name != "" {
		query += ", name = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Name)
	}
	if updatedCharacter.Affinity != "" {
		query += ", affinity = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Affinity)
	}
	if updatedCharacter.ImageLink != "" {
		query += ", image_link = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.ImageLink)
	}
	if updatedCharacter.BaseLv != 0 {
		query += ", base_lv = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.BaseLv)
	}
	if updatedCharacter.HP != 0 {
		query += ", hp = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.HP)
	}
	if updatedCharacter.HpGrowth != 0 {
		query += ", hp_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.HpGrowth)
	}
	if updatedCharacter.Strength != 0 {
		query += ", strength = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Strength)
	}
	if updatedCharacter.StrGrowth != 0 {
		query += ", str_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.StrGrowth)
	}
	if updatedCharacter.Magic != 0 {
		query += ", magic = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Magic)
	}
	if updatedCharacter.MagGrowth != 0 {
		query += ", mag_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.MagGrowth)
	}
	if updatedCharacter.Dexterity != 0 {
		query += ", dexterity = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Dexterity)
	}
	if updatedCharacter.DexGrowth != 0 {
		query += ", dex_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.DexGrowth)
	}
	if updatedCharacter.Speed != 0 {
		query += ", speed = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Speed)
	}
	if updatedCharacter.SpdGrowth != 0 {
		query += ", spd_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.SpdGrowth)
	}
	if updatedCharacter.Luck != 0 {
		query += ", luck = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Luck)
	}
	if updatedCharacter.LckGrowth != 0 {
		query += ", lck_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.LckGrowth)
	}
	if updatedCharacter.Defence != 0 {
		query += ", defence = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Defence)
	}
	if updatedCharacter.DefGrowth != 0 {
		query += ", def_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.DefGrowth)
	}
	if updatedCharacter.Resistance != 0 {
		query += ", resistance = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Resistance)
	}
	if updatedCharacter.ResGrowth != 0 {
		query += ", res_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.ResGrowth)
	}
	if updatedCharacter.Charm != 0 {
		query += ", charm = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.Charm)
	}
	if updatedCharacter.ChaGrowth != 0 {
		query += ", cha_growth = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedCharacter.ChaGrowth)
	}

	// Finish the query with the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)
	err := cc.db.QueryRow(query+" RETURNING *", args...).Scan(
		&updatedCharacter.ID, &updatedCharacter.Name, &updatedCharacter.ImageLink,
		&updatedCharacter.Affinity, &updatedCharacter.BaseLv, &updatedCharacter.HP,
		&updatedCharacter.HpGrowth, &updatedCharacter.Strength, &updatedCharacter.StrGrowth,
		&updatedCharacter.Magic, &updatedCharacter.MagGrowth, &updatedCharacter.Dexterity,
		&updatedCharacter.DexGrowth, &updatedCharacter.Speed, &updatedCharacter.SpdGrowth,
		&updatedCharacter.Luck, &updatedCharacter.LckGrowth, &updatedCharacter.Defence,
		&updatedCharacter.DefGrowth, &updatedCharacter.Resistance, &updatedCharacter.ResGrowth,
		&updatedCharacter.Charm, &updatedCharacter.ChaGrowth, &updatedCharacter.CreatedAt,
		&updatedCharacter.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cc *CharacterController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	charID := mux.Vars(r)["charID"]
	id, err := strconv.Atoi(charID)
	if err != nil {
		http.Error(w, "invalid character ID", http.StatusBadRequest)
		return
	}

	err = cc.deleteCharacter(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("Character with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting character: %s", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "msg": "Character deleted successfully."}`))
}

func (cc *CharacterController) deleteCharacter(id int) error {
	// Check if the character exists
	var exists bool
	err := cc.db.QueryRow("SELECT EXISTS (SELECT 1 FROM characters WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("character with ID %d not found", id)
	}

	// Character exists, proceed with deletion
	_, err = cc.db.Exec("DELETE FROM characters WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
