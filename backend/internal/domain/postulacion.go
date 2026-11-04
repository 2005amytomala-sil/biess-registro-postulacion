package domain

import "time"

type Postulacion struct {
	ID          int       `json:"id"`
	PostulanteID int      `json:"postulante_id"`
	Codigo      string    `json:"codigo"`
	Sector      string    `json:"sector"`
	Estado      string    `json:"estado"`
	MotivoCierre string   `json:"motivo_cierre,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DatosPersonales struct {
	PostulacionID   int       `json:"postulacion_id"`
	Nombres         string    `json:"nombres"`
	Apellidos       string    `json:"apellidos"`
	Cedula          string    `json:"cedula"`
	Genero          string    `json:"genero"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
	Nacionalidad    string    `json:"nacionalidad"`
	EstadoCivil     string    `json:"estado_civil"`
	ParejaNombres   string    `json:"pareja_nombres,omitempty"`
	ParejaCedula    string    `json:"pareja_cedula,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type DomicilioContacto struct {
	PostulacionID   int       `json:"postulacion_id"`
	PaisID          int       `json:"pais_id"`
	ProvinciaID     *int      `json:"provincia_id,omitempty"`
	CantonID        *int      `json:"canton_id,omitempty"`
	ProvinciaTexto  string    `json:"provincia_texto,omitempty"`
	CiudadTexto     string    `json:"ciudad_texto,omitempty"`
	Celular         string    `json:"celular"`
	Correo1         string    `json:"correo_1"`
	Correo2         string    `json:"correo_2,omitempty"`
	CallePrincipal  string    `json:"calle_principal,omitempty"`
	Numero          string    `json:"numero,omitempty"`
	CalleSecundaria string    `json:"calle_secundaria,omitempty"`
	Sector          string    `json:"sector,omitempty"`
	Referencia      string    `json:"referencia,omitempty"`
	TelefonoDomicilio string `json:"telefono_domicilio,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type FamiliarIESS struct {
	ID            int       `json:"id"`
	PostulacionID int       `json:"postulacion_id"`
	Nombres       string    `json:"nombres"`
	Parentesco    string    `json:"parentesco"`
	Institucion   string    `json:"institucion"`
	Area          string    `json:"area,omitempty"`
	Cargo         string    `json:"cargo,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type RespaldoOrganizacion struct {
	PostulacionID      int       `json:"postulacion_id"`
	TieneRespaldo      bool      `json:"tiene_respaldo"`
	NombreOrganizacion string    `json:"nombre_organizacion,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

type FormacionAcademica struct {
	ID                   int       `json:"id"`
	PostulacionID        int       `json:"postulacion_id"`
	Nivel                string    `json:"nivel"`
	Institucion          string    `json:"institucion"`
	Tipo                 string    `json:"tipo"`
	Titulo               string    `json:"titulo"`
	AreaEstudios         string    `json:"area_estudios,omitempty"`
	RegistroSenescyt     string    `json:"registro_senescyt,omitempty"`
	FechaRegistroSenescyt time.Time `json:"fecha_registro_senescyt,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

type ExperienciaProfesional struct {
	ID                     int       `json:"id"`
	PostulacionID          int       `json:"postulacion_id"`
	Tipo                   string    `json:"tipo"`
	Cargo                  string    `json:"cargo"`
	Institucion            string    `json:"institucion"`
	FechaInicio            time.Time `json:"fecha_inicio"`
	FechaFin               *time.Time `json:"fecha_fin,omitempty"`
	ActividadesRelacionadas string    `json:"actividades_relacionadas,omitempty"`
	Descripcion            string    `json:"descripcion,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
}

type Inhabilidad struct {
	PostulacionID  int       `json:"postulacion_id"`
	CodigoPregunta string    `json:"codigo_pregunta"`
	Respuesta      bool      `json:"respuesta"`
	Detalle        string    `json:"detalle,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type GuardarPostulacionRequest struct {
	Sector                string                    `json:"sector"`
	Domicilio             DomicilioContacto         `json:"domicilio"`
	TieneFamiliares       bool                      `json:"tiene_familiares"`
	Familiares            []FamiliarIESS            `json:"familiares,omitempty"`
	RespaldoOrganizacion  RespaldoOrganizacion     `json:"respaldo_organizacion"`
	Formacion             []FormacionAcademica       `json:"formacion,omitempty"`
	Experiencia           []ExperienciaProfesional   `json:"experiencia,omitempty"`
	Inhabilidades         []Inhabilidad             `json:"inhabilidades,omitempty"`
}
