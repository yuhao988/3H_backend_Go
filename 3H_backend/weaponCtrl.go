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

type WeaponsController struct {
	db *sql.DB
}

func NewWeaponsController(db *sql.DB) *WeaponsController {
	return &WeaponsController{
		db: db,
	}
}

func (cc *WeaponsController) GetAll(w http.ResponseWriter, r *http.Request) {
	weapons, err := cc.getAllWeapons()
	if err != nil {
		log.Printf("Error querying all weapons: %s", err)
		http.Error(w, fmt.Sprintf("Error getting weapons: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(weapons)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding weapons to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *WeaponsController) getAllWeapons() ([]Weapons, error) {

	rows, err := cc.db.Query("SELECT * FROM weapons")
	if err != nil {
		log.Printf("Error querying all weapons: %s", err)
		return nil, err
	}
	defer rows.Close()

	var weapons []Weapons

	for rows.Next() {
		var weapon Weapons

		err := rows.Scan(&weapon.ID, &weapon.Name, &weapon.TypeID,
			&weapon.StrMag, &weapon.Might, &weapon.Hit, &weapon.Critical,
			&weapon.Durability, &weapon.Weight, &weapon.RangeMin,
			&weapon.RangeMax, &weapon.Description, &weapon.CreatedAt,
			&weapon.UpdatedAt)
		if err != nil {
			return nil, err
		}
		weapons = append(weapons, weapon)
	}

	return weapons, nil
}

func (cc *WeaponsController) GetOne(w http.ResponseWriter, r *http.Request) {
	weaponID := mux.Vars(r)["weaponID"]
	id, err := strconv.Atoi(weaponID)
	if err != nil {
		http.Error(w, "Invalid weapon ID", http.StatusBadRequest)
		return
	}

	weapon, err := cc.getWeaponByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting weapon: %s", err), http.StatusInternalServerError)
		return
	}

	if weapon == nil {
		http.Error(w, "Weapon not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(weapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding weapon to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *WeaponsController) getWeaponByID(id int) (*Weapons, error) {
	// Implement the logic to fetch a weapon by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM weapons WHERE id = $1", id)

	var weapon Weapons
	err := row.Scan(&weapon.ID, &weapon.Name, &weapon.TypeID, &weapon.StrMag, &weapon.Might, &weapon.Hit, &weapon.Critical, &weapon.Durability, &weapon.Weight, &weapon.RangeMin, &weapon.RangeMax, &weapon.Description, &weapon.CreatedAt, &weapon.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &weapon, nil
}

func (cc *WeaponsController) GetOneName(w http.ResponseWriter, r *http.Request) {
	weaponName := mux.Vars(r)["weaponName"]

	weapon, err := cc.getWeaponByName(weaponName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting weapon: %s", err), http.StatusInternalServerError)
		return
	}

	if weapon == nil {
		http.Error(w, "Weapon not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(weapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding weapon to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *WeaponsController) getWeaponByName(name string) ([]Weapons, error) {
	// Implement the logic to fetch a weapon by ID from the database
	// Example:
	rows, err := cc.db.Query("SELECT * FROM weapons WHERE name LIKE $1", name+"%")
	if err != nil {
		log.Printf("Error querying all weapons: %s", err)
		return nil, err
	}
	defer rows.Close()

	var weapons []Weapons

	for rows.Next() {
		var weapon Weapons

		err := rows.Scan(&weapon.ID, &weapon.Name, &weapon.TypeID,
			&weapon.StrMag, &weapon.Might, &weapon.Hit, &weapon.Critical,
			&weapon.Durability, &weapon.Weight, &weapon.RangeMin,
			&weapon.RangeMax, &weapon.Description, &weapon.CreatedAt,
			&weapon.UpdatedAt)
		if err != nil {
			return nil, err
		}
		weapons = append(weapons, weapon)
	}

	return weapons, nil
}

func (cc *WeaponsController) PostOne(w http.ResponseWriter, r *http.Request) {
	var weapon Weapons
	err := json.NewDecoder(r.Body).Decode(&weapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	weapon.CreatedAt = time.Now()
	weapon.UpdatedAt = time.Now()

	err = cc.insertWeapon(&weapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting weapon: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(weapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding weapon to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *WeaponsController) insertWeapon(weapon *Weapons) error {
	// Perform the insert operation with the RETURNING clause to get the ID
	err := cc.db.QueryRow(`
		INSERT INTO weapons (name, type_id, str_mag, might, hit, critical, durability, 
			weight, range_min, range_max, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`, weapon.Name, weapon.TypeID, weapon.StrMag, weapon.Might, weapon.Hit,
		weapon.Critical, weapon.Durability, weapon.Weight, weapon.RangeMin, weapon.RangeMax,
		weapon.Description, weapon.CreatedAt, weapon.UpdatedAt).Scan(&weapon.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *WeaponsController) PutOne(w http.ResponseWriter, r *http.Request) {
	weaponID := mux.Vars(r)["weaponID"]
	id, err := strconv.Atoi(weaponID)
	if err != nil {
		http.Error(w, "Invalid weapon ID", http.StatusBadRequest)
		return
	}

	var updatedWeapon Weapons
	err = json.NewDecoder(r.Body).Decode(&updatedWeapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	updatedWeapon.UpdatedAt = time.Now()

	err = cc.updateWeapon(id, &updatedWeapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating weapon: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(updatedWeapon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated weapon to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *WeaponsController) updateWeapon(id int, updatedWeapon *Weapons) error {
	// Start building the SQL query
	query := "UPDATE weapons SET updated_at = $1"
	args := []interface{}{updatedWeapon.UpdatedAt}

	// Conditionally include fields in the update query
	if updatedWeapon.Name != "" {
		query += ", name = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.Name)
	}
	if updatedWeapon.TypeID != 0 {
		query += ", type_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.TypeID)
	}
	if updatedWeapon.StrMag != nil {
		query += ", str_mag = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.StrMag)
	}
	if updatedWeapon.Might != nil {
		query += ", might = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.Might)
	}
	if updatedWeapon.Hit != nil {
		query += ", hit = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.Hit)
	}
	if updatedWeapon.Critical != nil {
		query += ", critical = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.Critical)
	}
	if updatedWeapon.Durability != 0 {
		query += ", durability = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.Durability)
	}
	if updatedWeapon.Weight != 0 {
		query += ", weight = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.Weight)
	}
	if updatedWeapon.RangeMin != 0 {
		query += ", range_min = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.RangeMin)
	}
	if updatedWeapon.RangeMax != nil {
		query += ", range_max = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.RangeMax)
	}
	if updatedWeapon.Description != nil {
		query += ", description = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedWeapon.Description)
	}

	// Finish the query with the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)
	err := cc.db.QueryRow(query+" RETURNING *", args...).Scan(
		&updatedWeapon.ID, &updatedWeapon.Name, &updatedWeapon.TypeID, &updatedWeapon.StrMag,
		&updatedWeapon.Might, &updatedWeapon.Hit, &updatedWeapon.Critical,
		&updatedWeapon.Durability, &updatedWeapon.Weight, &updatedWeapon.RangeMin,
		&updatedWeapon.RangeMax, &updatedWeapon.Description, &updatedWeapon.CreatedAt,
		&updatedWeapon.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cc *WeaponsController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	weaponID := mux.Vars(r)["weaponID"]
	id, err := strconv.Atoi(weaponID)
	if err != nil {
		http.Error(w, "Invalid weapon ID", http.StatusBadRequest)
		return
	}

	err = cc.deleteWeapon(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("Weapon with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting weapon: %s", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "msg": "Weapon deleted successfully."}`))
}

func (cc *WeaponsController) deleteWeapon(id int) error {
	// Check if the weapon exists
	var exists bool
	err := cc.db.QueryRow("SELECT EXISTS (SELECT 1 FROM weapons WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("weapon with ID %d not found", id)
	}

	// Weapon exists, proceed with deletion
	_, err = cc.db.Exec("DELETE FROM weapons WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
