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

type SkillsController struct {
	db *sql.DB
}

func NewSkillsController(db *sql.DB) *SkillsController {
	return &SkillsController{
		db: db,
	}
}

func (cc *SkillsController) GetAll(w http.ResponseWriter, r *http.Request) {
	skills, err := cc.getAllSkills()
	if err != nil {
		log.Printf("Error querying all skill types: %s", err)
		http.Error(w, fmt.Sprintf("Error getting skill types: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(skills)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding skill types to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SkillsController) getAllSkills() ([]Skills, error) {

	rows, err := cc.db.Query("SELECT * FROM skills")
	if err != nil {
		log.Printf("Error querying all skill types: %s", err)
		return nil, err
	}
	defer rows.Close()

	var skills []Skills

	for rows.Next() {
		var skill Skills
		err := rows.Scan(&skill.ID, &skill.Name, &skill.SkillIcon, &skill.CreatedAt, &skill.UpdatedAt)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (cc *SkillsController) GetOne(w http.ResponseWriter, r *http.Request) {
	skillID := mux.Vars(r)["skillID"]
	id, err := strconv.Atoi(skillID)
	if err != nil {
		http.Error(w, "Invalid skill type ID", http.StatusBadRequest)
		return
	}

	skill, err := cc.getSkillByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting skill type: %s", err), http.StatusInternalServerError)
		return
	}

	if skill == nil {
		http.Error(w, "Skill type not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(skill)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding skill type  to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SkillsController) getSkillByID(id int) (*Skills, error) {
	// Implement the logic to fetch a skill type  by ID from the database
	// Example:
	row := cc.db.QueryRow("SELECT * FROM skills WHERE id = $1", id)

	var skill Skills
	err := row.Scan(&skill.ID, &skill.Name, &skill.SkillIcon, &skill.CreatedAt, &skill.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &skill, nil
}

func (cc *SkillsController) PostOne(w http.ResponseWriter, r *http.Request) {
	var skill Skills
	err := json.NewDecoder(r.Body).Decode(&skill)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	skill.CreatedAt = time.Now()
	skill.UpdatedAt = time.Now()

	err = cc.insertSkill(&skill)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting skill type : %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(skill)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding skill type  to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SkillsController) insertSkill(skill *Skills) error {
	// Perform the insert operation with the RETURNING clause to get the ID
	err := cc.db.QueryRow(`
		INSERT INTO skills (name, skill_icon, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, skill.Name, skill.SkillIcon, skill.CreatedAt, skill.UpdatedAt).Scan(&skill.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *SkillsController) PutOne(w http.ResponseWriter, r *http.Request) {
	skillID := mux.Vars(r)["skillID"]
	id, err := strconv.Atoi(skillID)
	if err != nil {
		http.Error(w, "Invalid skill type ID", http.StatusBadRequest)
		return
	}

	var updatedSkill Skills
	err = json.NewDecoder(r.Body).Decode(&updatedSkill)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	updatedSkill.UpdatedAt = time.Now()

	err = cc.updateSkill(id, &updatedSkill)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating skill type: %s", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(updatedSkill)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated skill type to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (cc *SkillsController) updateSkill(id int, updatedSkill *Skills) error {
	// Start building the SQL query
	query := "UPDATE skills SET updated_at = $1"
	args := []interface{}{updatedSkill.UpdatedAt}

	// Conditionally include fields in the update query
	if updatedSkill.Name != "" {
		query += ", name = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSkill.Name)
	}
	if updatedSkill.SkillIcon != nil {
		query += ", skill_icon = $" + strconv.Itoa(len(args)+1)
		args = append(args, updatedSkill.SkillIcon)
	}

	// Finish the query with the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)
	err := cc.db.QueryRow(query+" RETURNING *", args...).Scan(
		&updatedSkill.ID, &updatedSkill.Name,
		&updatedSkill.SkillIcon, &updatedSkill.CreatedAt,
		&updatedSkill.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cc *SkillsController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	skillID := mux.Vars(r)["skillID"]
	id, err := strconv.Atoi(skillID)
	if err != nil {
		http.Error(w, "Invalid skill type ID", http.StatusBadRequest)
		return
	}

	err = cc.deleteSkill(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("Skill type with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting skill type: %s", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "msg": "Skill type deleted successfully."}`))
}

func (cc *SkillsController) deleteSkill(id int) error {
	// Check if the skill type  exists
	var exists bool
	err := cc.db.QueryRow("SELECT EXISTS (SELECT 1 FROM skills WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("skill type with ID %d not found", id)
	}

	// Skill exists, proceed with deletion
	_, err = cc.db.Exec("DELETE FROM skills WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
