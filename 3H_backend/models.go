package main

import (
	"time"
)

type Character struct {
	ID   int
	Name string // Update with your actual fields

	ImageLink  string
	Affinity   string
	BaseLv     int
	HP         int
	HpGrowth   int
	Strength   int
	StrGrowth  int
	Magic      int
	MagGrowth  int
	Dexterity  int
	DexGrowth  int
	Speed      int
	SpdGrowth  int
	Luck       int
	LckGrowth  int
	Defence    int
	DefGrowth  int
	Resistance int
	ResGrowth  int
	Charm      int
	ChaGrowth  int
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Spells struct {
	ID   int
	Name string
	Type string
	// * before the type indicates that the value is nullable
	Might       *int
	Hit         *int
	Critical    *int
	Uses        int
	Weight      *int
	RangeMin    int
	RangeMax    *int
	Description *string
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Skills struct {
	ID        int
	Name      string
	SkillIcon *string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Secondary
type CombatArt struct{}

type Weapon struct{}

// Tertiary
type CharSkill struct{}
