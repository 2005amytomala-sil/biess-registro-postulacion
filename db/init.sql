-- ============================================================================
-- BLOQUE 1: DDL Normalizado (3FN)
-- Sistema de Registro y Postulación al Concurso Público de Méritos del BIESS
-- ============================================================================

-- Habilitar extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- CATÁLOGOS: Paises, Provincias, Cantones
-- ============================================================================

CREATE TABLE paises (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE provincias (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cantones (
    id SERIAL PRIMARY KEY,
    provincia_id INTEGER NOT NULL REFERENCES provincias(id) ON DELETE RESTRICT,
    nombre VARCHAR(100) NOT NULL,
    UNIQUE (provincia_id, nombre),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para catálogos geográficos
CREATE INDEX idx_cantones_provincia_id ON cantones(provincia_id);

-- ============================================================================
-- TABLA: postulantes
-- ============================================================================
CREATE TABLE postulantes (
    id SERIAL PRIMARY KEY,
    cedula VARCHAR(20) NOT NULL UNIQUE,
    correo VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para postulantes
CREATE INDEX idx_postulantes_cedula ON postulantes(cedula);
CREATE INDEX idx_postulantes_correo ON postulantes(correo);

-- ============================================================================
-- TABLA: postulaciones
-- ============================================================================
CREATE TABLE postulaciones (
    id SERIAL PRIMARY KEY,
    postulante_id INTEGER NOT NULL UNIQUE REFERENCES postulantes(id) ON DELETE CASCADE,
    codigo VARCHAR(50) NOT NULL UNIQUE,
    sector VARCHAR(20) NOT NULL CHECK (sector IN ('afiliados', 'jubilados', 'empleadores')),
    estado VARCHAR(20) NOT NULL DEFAULT 'borrador' CHECK (estado IN ('borrador', 'enviada', 'cerrada')),
    motivo_cierre TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para postulaciones
CREATE INDEX idx_postulaciones_postulante_id ON postulaciones(postulante_id);
CREATE INDEX idx_postulaciones_codigo ON postulaciones(codigo);
CREATE INDEX idx_postulaciones_estado ON postulaciones(estado);

-- Trigger para actualizar updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_postulaciones_updated_at
    BEFORE UPDATE ON postulaciones
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- TABLA: datos_personales
-- ============================================================================
CREATE TABLE datos_personales (
    postulacion_id INTEGER PRIMARY KEY REFERENCES postulaciones(id) ON DELETE CASCADE,
    nombres VARCHAR(200) NOT NULL,
    apellidos VARCHAR(200) NOT NULL,
    cedula VARCHAR(20) NOT NULL,
    genero VARCHAR(20) NOT NULL CHECK (genero IN ('masculino', 'femenino', 'otro')),
    fecha_nacimiento DATE NOT NULL,
    nacionalidad VARCHAR(50) NOT NULL,
    estado_civil VARCHAR(50) NOT NULL,
    pareja_nombres VARCHAR(200),
    pareja_cedula VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- TABLA: domicilio_contacto
-- ============================================================================
CREATE TABLE domicilio_contacto (
    postulacion_id INTEGER PRIMARY KEY REFERENCES postulaciones(id) ON DELETE CASCADE,
    pais_id INTEGER NOT NULL REFERENCES paises(id) ON DELETE RESTRICT,
    provincia_id INTEGER REFERENCES provincias(id) ON DELETE RESTRICT,
    canton_id INTEGER REFERENCES cantones(id) ON DELETE RESTRICT,
    provincia_texto VARCHAR(100),
    ciudad_texto VARCHAR(100),
    celular VARCHAR(20) NOT NULL,
    correo_1 VARCHAR(255) NOT NULL,
    correo_2 VARCHAR(255),
    calle_principal VARCHAR(200),
    numero VARCHAR(50),
    calle_secundaria VARCHAR(200),
    sector VARCHAR(200),
    referencia TEXT,
    telefono_domicilio VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (pais_id = 1 AND provincia_id IS NOT NULL AND canton_id IS NOT NULL) OR
        (pais_id != 1 AND provincia_id IS NULL AND canton_id IS NULL)
    )
);

-- Índices para domicilio_contacto
CREATE INDEX idx_domicilio_pais_id ON domicilio_contacto(pais_id);
CREATE INDEX idx_domicilio_provincia_id ON domicilio_contacto(provincia_id);
CREATE INDEX idx_domicilio_canton_id ON domicilio_contacto(canton_id);

-- ============================================================================
-- TABLA: familiares_iess
-- ============================================================================
CREATE TABLE familiares_iess (
    id SERIAL PRIMARY KEY,
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    nombres VARCHAR(200) NOT NULL,
    parentesco VARCHAR(100) NOT NULL,
    institucion VARCHAR(20) NOT NULL CHECK (institucion IN ('IESS', 'BIESS')),
    area VARCHAR(200),
    cargo VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para familiares_iess
CREATE INDEX idx_familiares_postulacion_id ON familiares_iess(postulacion_id);

-- ============================================================================
-- TABLA: respaldo_organizacion
-- ============================================================================
CREATE TABLE respaldo_organizacion (
    postulacion_id INTEGER PRIMARY KEY REFERENCES postulaciones(id) ON DELETE CASCADE,
    tiene_respaldo BOOLEAN NOT NULL DEFAULT FALSE,
    nombre_organizacion VARCHAR(300),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (tiene_respaldo = TRUE AND nombre_organizacion IS NOT NULL) OR
        (tiene_respaldo = FALSE)
    )
);

-- ============================================================================
-- TABLA: formacion_academica
-- ============================================================================
CREATE TABLE formacion_academica (
    id SERIAL PRIMARY KEY,
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    nivel VARCHAR(100) NOT NULL,
    institucion VARCHAR(300) NOT NULL,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('Nacional', 'Extranjera')),
    titulo VARCHAR(300) NOT NULL,
    area_estudios VARCHAR(200),
    registro_senescyt VARCHAR(100),
    fecha_registro_senescyt DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para formacion_academica
CREATE INDEX idx_formacion_postulacion_id ON formacion_academica(postulacion_id);

-- ============================================================================
-- TABLA: experiencia_profesional
-- ============================================================================
CREATE TABLE experiencia_profesional (
    id SERIAL PRIMARY KEY,
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('General', 'Direccion')),
    cargo VARCHAR(200) NOT NULL,
    institucion VARCHAR(300) NOT NULL,
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE,
    actividades_relacionadas TEXT,
    descripcion VARCHAR(1000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (fecha_fin IS NULL OR fecha_fin >= fecha_inicio)
);

-- Índices para experiencia_profesional
CREATE INDEX idx_experiencia_postulacion_id ON experiencia_profesional(postulacion_id);

-- ============================================================================
-- TABLA: inhabilidades
-- ============================================================================
CREATE TABLE inhabilidades (
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    codigo_pregunta VARCHAR(20) NOT NULL CHECK (codigo_pregunta IN (
        '1.9.1', '1.9.2', '1.9.3', '1.9.4', '1.9.5',
        '1.9.6', '1.9.7', '1.9.8', '1.9.9', '1.9.10',
        '1.9.11', '1.9.12', '1.9.13'
    )),
    respuesta BOOLEAN NOT NULL,
    detalle JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (postulacion_id, codigo_pregunta),
    CHECK (
        (codigo_pregunta = '1.9.7' AND detalle IS NOT NULL) OR
        (codigo_pregunta != '1.9.7')
    )
);

-- Índices para inhabilidades
CREATE INDEX idx_inhabilidades_postulacion_id ON inhabilidades(postulacion_id);
CREATE INDEX idx_inhabilidades_codigo_pregunta ON inhabilidades(codigo_pregunta);

-- ============================================================================
-- FIN DEL BLOQUE 1: DDL
-- ============================================================================
