package model

import (
	"time"
)

type Employee struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Name              string     `json:"name" gorm:"not null;column:name;size:100"`
	Email             string     `json:"email" gorm:"not null;column:email;size:100"`
	ContactNumber     string     `json:"contact_number" gorm:"column:contact_number;size:15"`
	Gender            string     `json:"gender" gorm:"column:gender;size:15"`
	Department        string     `json:"department" gorm:"column:department"`
	TechStack         string     `json:"tech_stack" gorm:"column:tech_stack;size:255"`
	DateOfJoining     *time.Time `json:"date_of_joining" gorm:"column:date_of_joining"`
	Position          string     `json:"position" gorm:"column:position;size:100"`
	YearsOfExperience float32    `json:"YearsOfExperience" gorm:"column:years_of_experience"`
	CasualLeave       int        `json:"CasualLeave" gorm:"column:casual_leave"`
	EarnedLeave       int        `json:"EarnedLeave" gorm:"column:earned_leave"`
	Salary            float64    `json:"salary" gorm:"column:salary"`
	Performance       string     `json:"performance" gorm:"column:performance;size:255"`
	BirthDate         *time.Time `json:"birth_date" gorm:"column:birth_date"`
	Address           string     `json:"address" gorm:"column:address;size:255"`
	Resume            string     `json:"resume" gorm:"column:resume"`
	ExperienceLetter  string     `json:"experience_letter" gorm:"column:experience_letter"`
	ReleivingLetter   string     `json:"releiving_letter" gorm:"column:releiving_letter"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
