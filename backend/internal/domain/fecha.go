package domain

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const formatoFecha = "2006-01-02"

// Date representa una fecha sin hora ("YYYY-MM-DD"). Acepta vacío o null como fecha ausente.
type Date struct {
	Time  time.Time
	Valid bool
}

var ErrFechaFormatoInvalido = NewAppError("FECHA_FORMATO_INVALIDO", "Formato de fecha inválido. Use AAAA-MM-DD", nil)

func ParseDate(s string) (Date, error) {
	if s == "" {
		return Date{}, nil
	}
	t, err := time.Parse(formatoFecha, s)
	if err != nil {
		return Date{}, ErrFechaFormatoInvalido
	}
	return Date{Time: t, Valid: true}, nil
}

func (d Date) String() string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format(formatoFecha)
}

// Before indica si d es anterior a o (ambas deben ser válidas)
func (d Date) Before(o Date) bool {
	return d.Valid && o.Valid && d.Time.Before(o.Time)
}

func (d *Date) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		*d = Date{}
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return ErrFechaFormatoInvalido
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.String())
}

// Value implementa driver.Valuer (NULL si no es válida)
func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time, nil
}

// Scan implementa sql.Scanner
func (d *Date) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*d = Date{}
	case time.Time:
		*d = Date{Time: time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC), Valid: true}
	case string:
		parsed, err := ParseDate(v)
		if err != nil {
			return err
		}
		*d = parsed
	case []byte:
		return d.Scan(string(v))
	default:
		return fmt.Errorf("no se puede convertir %T a Date", src)
	}
	return nil
}
