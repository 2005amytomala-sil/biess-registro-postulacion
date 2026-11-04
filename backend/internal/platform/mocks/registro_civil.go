package mocks

import (
	"context"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
)

type RegistroCivilMock struct {
	dataset map[string]*domain.RegistroCivilData
}

func NewRegistroCivilMock() *RegistroCivilMock {
	dataset := make(map[string]*domain.RegistroCivilData)

	// Caso 1: Válido - ecuatoriano, soltero
	dataset["0101010101"] = &domain.RegistroCivilData{
		Cedula:            "0101010101",
		FechaExpedicion:   "01/01/2010",
		CodigoDactilar:    "1234567890",
		Nombres:           "Juan Carlos",
		Apellidos:         "Pérez López",
		Genero:            "masculino",
		FechaNacimiento:   "01/01/1980",
		Nacionalidad:      "Ecuatoriana",
		EstadoCivil:       "soltero",
	}

	// Caso 2: Válido - ecuatoriana, soltera
	dataset["0202020202"] = &domain.RegistroCivilData{
		Cedula:            "0202020202",
		FechaExpedicion:   "15/06/2012",
		CodigoDactilar:    "0987654321",
		Nombres:           "María Isabel",
		Apellidos:         "García Martínez",
		Genero:            "femenino",
		FechaNacimiento:   "15/06/1985",
		Nacionalidad:      "Ecuatoriana",
		EstadoCivil:       "soltero",
	}

	// Caso 3: Nacionalidad no ecuatoriana
	dataset["0303030303"] = &domain.RegistroCivilData{
		Cedula:            "0303030303",
		FechaExpedicion:   "20/03/2015",
		CodigoDactilar:    "1122334455",
		Nombres:           "John",
		Apellidos:         "Smith",
		Genero:            "masculino",
		FechaNacimiento:   "20/03/1975",
		Nacionalidad:      "Estadounidense",
		EstadoCivil:       "soltero",
	}

	// Caso 4: Casado con datos de pareja
	dataset["0404040404"] = &domain.RegistroCivilData{
		Cedula:            "0404040404",
		FechaExpedicion:   "10/09/2018",
		CodigoDactilar:    "5566778899",
		Nombres:           "Carlos Eduardo",
		Apellidos:         "Rodríguez Torres",
		Genero:            "masculino",
		FechaNacimiento:   "10/09/1978",
		Nacionalidad:      "Ecuatoriana",
		EstadoCivil:       "casado",
		ParejaNombres:     "Ana María",
		ParejaCedula:      "0505050505",
	}

	// Caso 5: Casada con datos de pareja
	dataset["0505050505"] = &domain.RegistroCivilData{
		Cedula:            "0505050505",
		FechaExpedicion:   "05/12/2019",
		CodigoDactilar:    "9988776655",
		Nombres:           "Ana María",
		Apellidos:         "Gómez Sánchez",
		Genero:            "femenino",
		FechaNacimiento:   "05/12/1982",
		Nacionalidad:      "Ecuatoriana",
		EstadoCivil:       "casado",
		ParejaNombres:     "Carlos Eduardo",
		ParejaCedula:      "0404040404",
	}

	// Caso 6: Válido - ecuatoriano, divorciado
	dataset["0606060606"] = &domain.RegistroCivilData{
		Cedula:            "0606060606",
		FechaExpedicion:   "25/04/2020",
		CodigoDactilar:    "4455667788",
		Nombres:           "Luis Fernando",
		Apellidos:         "Mendoza Ruiz",
		Genero:            "masculino",
		FechaNacimiento:   "25/04/1979",
		Nacionalidad:      "Ecuatoriana",
		EstadoCivil:       "divorciado",
	}

	// Caso 7: Válido - ecuatoriano, soltero (para pruebas de bloqueo)
	dataset["0909090909"] = &domain.RegistroCivilData{
		Cedula:            "0909090909",
		FechaExpedicion:   "01/01/2015",
		CodigoDactilar:    "1111222233",
		Nombres:           "Carlos",
		Apellidos:         "Test Prueba",
		Genero:            "masculino",
		FechaNacimiento:   "01/01/1990",
		Nacionalidad:      "Ecuatoriana",
		EstadoCivil:       "soltero",
	}

	return &RegistroCivilMock{dataset: dataset}
}

func (m *RegistroCivilMock) VerificarIdentidad(ctx context.Context, cedula string, fechaExpedicion string, codigoDactilar string) (*domain.RegistroCivilData, error) {
	data, exists := m.dataset[cedula]
	if !exists {
		return nil, domain.ErrIdentidadInvalida
	}

	// Verificar fecha de expedición si se proporciona
	if fechaExpedicion != "" && data.FechaExpedicion != fechaExpedicion {
		return nil, domain.ErrIdentidadInvalida
	}

	// Verificar código dactilar si se proporciona
	if codigoDactilar != "" && data.CodigoDactilar != codigoDactilar {
		return nil, domain.ErrIdentidadInvalida
	}

	return data, nil
}

func (m *RegistroCivilMock) ObtenerDatos(ctx context.Context, cedula string) (*domain.RegistroCivilData, error) {
	data, exists := m.dataset[cedula]
	if !exists {
		return nil, domain.ErrIdentidadInvalida
	}
	return data, nil
}
