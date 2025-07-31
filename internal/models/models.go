package models

import "time"

type User struct {
	ID        uint `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          uint `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
	Stock       int `json:"stock"`
	Tags        []string `json:"tags"`
}
