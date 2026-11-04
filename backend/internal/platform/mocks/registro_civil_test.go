package mocks

import (
	"context"
	"testing"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
)

func TestRegistroCivilMock_ObtenerDatos(t *testing.T) {
	mock := NewRegistroCivilMock()

	ctx := context.Background()

	// Test con cédula existente
	data, err := mock.ObtenerDatos(ctx, "0101010101")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if data == nil {
		t.Error("Expected data to not be nil")
	}
	if data.Cedula != "0101010101" {
		t.Errorf("Expected cedula 0101010101, got %s", data.Cedula)
	}
	if data.Nombres != "Juan Carlos" {
		t.Errorf("Expected nombres Juan Carlos, got %s", data.Nombres)
	}

	// Test con cédula inexistente
	data, err = mock.ObtenerDatos(ctx, "9999999999")
	if err == nil {
		t.Error("Expected error for non-existent cedula")
	}
	if err != domain.ErrIdentidadInvalida {
		t.Errorf("Expected ErrIdentidadInvalida, got %v", err)
	}
	if data != nil {
		t.Error("Expected data to be nil for non-existent cedula")
	}
}
