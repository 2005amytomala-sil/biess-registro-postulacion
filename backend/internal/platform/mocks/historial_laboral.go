package mocks

import "context"

type HistorialLaboralClient interface {
	VerificarSector(ctx context.Context, cedula string, sector string) (bool, string, error)
}

type HistorialLaboralMock struct{}

func NewHistorialLaboralMock() *HistorialLaboralMock {
	return &HistorialLaboralMock{}
}

func (m *HistorialLaboralMock) VerificarSector(ctx context.Context, cedula string, sector string) (bool, string, error) {
	// Mock: Simula validación de sector
	// Para el test, asumimos que ciertas cédulas tienen ciertos sectores
	// Por defecto, devolvemos verdadero para afiliados
	
	switch sector {
	case "afiliados":
		// Simular que algunas cédulas no son afiliados
		if cedula == "0303030303" {
			return false, "Su postulación para el cargo de Miembro Principal y Alterno en Representación del Sector de Afiliados ante el Directorio del BIESS, NO puede continuar, por cuanto no mantiene la calidad de afiliado al IESS", nil
		}
		return true, "", nil
	case "jubilados":
		if cedula == "0404040404" {
			return false, "Su postulación para el cargo de Miembro Principal y Alterno en Representación del Sector de Jubilados ante el Directorio del BIESS, NO puede continuar, por cuanto no mantiene la calidad de jubilado del IESS", nil
		}
		return true, "", nil
	case "empleadores":
		if cedula == "0505050505" {
			return false, "Su postulación para el cargo de Miembro Principal y Alterno en Representación del Sector de Empleadores ante el Directorio del BIESS, NO puede continuar, por cuanto no mantiene la calidad de empleador afiliado al IESS", nil
		}
		return true, "", nil
	default:
		return false, "Sector no válido", nil
	}
}
