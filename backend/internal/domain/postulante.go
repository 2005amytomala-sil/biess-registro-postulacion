package domain

import "time"

type Postulante struct {
	ID           int       `json:"id"`
	Cedula       string    `json:"cedula"`
	Correo       string    `json:"correo"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegistroCivilData struct {
	Cedula            string `json:"cedula"`
	FechaExpedicion   string `json:"fecha_expedicion"`
	CodigoDactilar    string `json:"codigo_dactilar"`
	Nombres           string `json:"nombres"`
	Apellidos         string `json:"apellidos"`
	Genero            string `json:"genero"`
	FechaNacimiento   string `json:"fecha_nacimiento"`
	Nacionalidad      string `json:"nacionalidad"`
	EstadoCivil       string `json:"estado_civil"`
	ParejaNombres     string `json:"pareja_nombres,omitempty"`
	ParejaCedula      string `json:"pareja_cedula,omitempty"`
}

type VerificarIdentidadRequest struct {
	Cedula          string `json:"cedula" binding:"required"`
	FechaExpedicion string `json:"fecha_expedicion"`
	CodigoDactilar  string `json:"codigo_dactilar"`
}

type RegistroRequest struct {
	Cedula          string `json:"cedula" binding:"required"`
	Correo          string `json:"correo" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginRequest struct {
	Cedula    string `json:"cedula" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type RecuperarRequest struct {
	Cedula          string `json:"cedula" binding:"required"`
	FechaExpedicion string `json:"fecha_expedicion"`
	CodigoDactilar  string `json:"codigo_dactilar"`
}

type RestablecerRequest struct {
	Cedula          string `json:"cedula" binding:"required"`
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  Postulante   `json:"user"`
}

type BloqueoResponse struct {
	Mensaje      string `json:"mensaje"`
	RetryAfter   int    `json:"retry_after"`
}
