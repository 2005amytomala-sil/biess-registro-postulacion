package service

import "github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"

// Verificación en compilación de que los servicios implementan los casos de uso de dominio
var (
	_ domain.AuthService        = (*AuthService)(nil)
	_ domain.CatalogoService    = (*CatalogoService)(nil)
	_ domain.PostulacionService = (*PostulacionService)(nil)
)

type Service struct {
	// Dependencies will be injected here
}

func NewService() *Service {
	return &Service{}
}
