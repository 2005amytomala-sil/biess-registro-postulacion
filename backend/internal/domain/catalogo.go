package domain

import "time"

type Pais struct {
	ID        int       `json:"id"`
	Nombre    string    `json:"nombre"`
	CreatedAt time.Time `json:"created_at"`
}

type Provincia struct {
	ID        int       `json:"id"`
	Nombre    string    `json:"nombre"`
	CreatedAt time.Time `json:"created_at"`
}

type Canton struct {
	ID          int       `json:"id"`
	ProvinciaID int       `json:"provincia_id"`
	Nombre      string    `json:"nombre"`
	CreatedAt   time.Time `json:"created_at"`
}
