package model

import (
	"time"
)

// User Struct
type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	MobilePhone string    `json:"mobile_phone"`
	Password    string    `json:"password"`
	FullName    string    `json:"full_name"`
	Name        string    `json:"name"`
	Dob         time.Time `json:"dob"`
	Photo       string    `json:"photo"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IdDevice    string    `json:"id_device"`
	Pin         string    `json:"pin"`
}

// User type User list
type Users []User

func NewUser() *User {
	return &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
