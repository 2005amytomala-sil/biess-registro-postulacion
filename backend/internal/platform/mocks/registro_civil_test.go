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

func TestRegistroCivilMock_VerificarIdentidad(t *testing.T) {
	m := NewRegistroCivilMock()
	ctx := context.Background()

	casos := []struct {
		nombre, cedula, fecha, codigo string
		ok                            bool
	}{
		{"ambos vacios", "0101010101", "", "", false},
		{"fecha correcta", "0101010101", "01/01/2010", "", true},
		{"codigo correcto", "0101010101", "", "1234567890", true},
		{"ambos correctos", "0101010101", "01/01/2010", "1234567890", true},
		{"fecha incorrecta", "0101010101", "02/01/2010", "", false},
		{"codigo incorrecto", "0101010101", "", "0000000000", false},
		{"fecha ok codigo incorrecto", "0101010101", "01/01/2010", "0000000000", false},
		{"cedula inexistente", "9999999999", "01/01/2010", "", false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			data, err := m.VerificarIdentidad(ctx, c.cedula, c.fecha, c.codigo)
			if c.ok && (err != nil || data == nil) {
				t.Fatalf("esperaba éxito, got err=%v", err)
			}
			if !c.ok && err != domain.ErrIdentidadInvalida {
				t.Fatalf("esperaba ErrIdentidadInvalida, got %v", err)
			}
		})
	}
}
